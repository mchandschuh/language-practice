#!/usr/bin/env bash
# Go bootstrap for the language-practice track.
#
# This is a visible, user-run guide. The iterator and reviewer must inspect and
# statically validate this file only; they must not run it, install Go, invoke
# sudo, change PATH, initialize a module, or resolve/download Go modules.
# Run it explicitly from the repository root after reading every command.
#
# Official references:
# - Go installation: https://go.dev/doc/install
# - Go modules and dependencies: https://go.dev/doc/modules/managing-dependencies
# - Go diagnostics: https://go.dev/doc/diagnostics
#
# The repository has no selected Go support version. Go 1.27 is a documentation
# observation only, not an execution contract. Choose and record the version,
# architecture, source, permissions, and checksum for the method you run.
#
# The default invocation installs a missing Go host toolchain with an explicit
# detected distribution-package command. Archive installation is opt-in and
# requires both a version and a SHA-256 checksum. No shell profile is edited.

set -Eeuo pipefail

readonly TRACK_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
readonly REPOSITORY_ROOT="$(cd -- "$TRACK_DIR/.." && pwd -P)"

GO_BIN=''
INSTALL_METHOD='auto'
GO_VERSION=''
ARCHIVE_SHA256=''
INSTALL_DIR='/usr/local/go'
ALLOW_REPLACE=0
NO_INSTALL=0
MODULE_PATH=''
RUN_TIDY=0
RUN_TEST=0
RUN_VET=0
GO_INSTALL_STATUS='already present'
ARCHIVE_TEMP_DIR=''

die() {
  printf 'bootstrap error: %s\n' "$*" >&2
  exit 1
}

note() {
  printf '[go/bootstrap] %s\n' "$*"
}

cleanup_archive() {
  if [[ -n "$ARCHIVE_TEMP_DIR" && -d "$ARCHIVE_TEMP_DIR" ]]; then
    rm -rf -- "$ARCHIVE_TEMP_DIR"
  fi
}

trap cleanup_archive EXIT

usage() {
  cat <<'EOF'
Usage: ./go/bootstrap.sh [options]

Installs a missing Go toolchain, then performs explicitly requested project-local
module and verification actions. Run from the repository root; project state is
kept under the ./go/ track directory.

Toolchain options:
  --go PATH              Go executable or command name to use
  --install-method NAME  auto, apt, dnf, pacman, archive, or none (default: auto)
  --version VERSION      Go archive version, required with archive mode
  --sha256 HEX           64-character archive checksum, required with archive mode
  --install-dir PATH     archive destination (default: /usr/local/go)
  --replace              allow archive mode to replace only --install-dir
  --no-install           fail instead of installing a missing Go toolchain

Project options:
  --module PATH          explicitly run `go mod init PATH` when go.mod is absent
  --tidy                 run `go mod tidy` and permit module resolution/downloads
  --test                 run bounded `go test ./...`
  --vet                  run bounded `go vet ./...`
  --help                 show this help

Installation methods:
  auto       Detect apt, dnf, or pacman and run its visible Go package command.
  apt        Run `sudo apt-get update` and `sudo apt-get install -y golang-go`.
  dnf        Run `sudo dnf install -y golang`.
  pacman     Run `sudo pacman -S --needed --noconfirm go`.
  archive    Download an official archive, verify --sha256, and install it.
  none       Never install; fail if Go is unavailable.

The script never edits PATH or shell profiles. After a user-local archive
installation, invoke the reported Go path directly or add it to PATH manually
in the current shell after reviewing the command.
EOF
}

resolve_command() {
  local requested=$1
  if [[ "$requested" == */* ]]; then
    [[ -x "$requested" ]] || die "Go executable is not executable: $requested"
    printf '%s\n' "$requested"
  else
    command -v -- "$requested" || die "Go command not found: $requested"
  fi
}

detect_install_method() {
  case "$INSTALL_METHOD" in
    auto)
      if command -v apt-get >/dev/null 2>&1; then
        INSTALL_METHOD='apt'
      elif command -v dnf >/dev/null 2>&1; then
        INSTALL_METHOD='dnf'
      elif command -v pacman >/dev/null 2>&1; then
        INSTALL_METHOD='pacman'
      else
        die 'cannot detect apt-get, dnf, or pacman; choose --install-method archive with --version and --sha256'
      fi
      ;;
    apt|dnf|pacman|archive|none) ;;
    *) die "unsupported --install-method: $INSTALL_METHOD" ;;
  esac
}

require_sudo() {
  command -v sudo >/dev/null 2>&1 || die "sudo is required for $INSTALL_METHOD installation"
}

install_from_package() {
  require_sudo
  case "$INSTALL_METHOD" in
    apt)
      note 'installing Go host prerequisites: sudo apt-get update'
      sudo apt-get update
      note 'installing Go toolchain: sudo apt-get install -y golang-go'
      sudo apt-get install -y golang-go
      ;;
    dnf)
      note 'installing Go toolchain: sudo dnf install -y golang'
      sudo dnf install -y golang
      ;;
    pacman)
      note 'installing Go toolchain: sudo pacman -S --needed --noconfirm go'
      sudo pacman -S --needed --noconfirm go
      ;;
    *) die "internal package installation error for method: $INSTALL_METHOD" ;;
  esac
  GO_INSTALL_STATUS='installed by this run'
}

archive_architecture() {
  case "$(uname -m)" in
    x86_64) printf '%s\n' 'amd64' ;;
    aarch64|arm64) printf '%s\n' 'arm64' ;;
    armv6l|armv7l) printf '%s\n' 'armv6l' ;;
    *) die "unsupported Linux architecture for Go archive: $(uname -m)" ;;
  esac
}

validate_archive_destination() {
  [[ "$INSTALL_DIR" == /* ]] || die '--install-dir must be an absolute path'
  [[ "$INSTALL_DIR" != */ ]] || die '--install-dir must not end with /'
  case "$INSTALL_DIR" in
    /|/usr|/usr/local|/opt|/bin|/sbin|/etc|/var|/home|/root)
      die "refusing to use a system root as --install-dir: $INSTALL_DIR"
      ;;
  esac
  case "$INSTALL_DIR" in
    *'//'*) die "--install-dir must not contain an empty path component: $INSTALL_DIR" ;;
    */..|*/../*) die "--install-dir must not contain .. path components: $INSTALL_DIR" ;;
  esac
}

install_from_archive() {
  [[ "$GO_VERSION" =~ ^[0-9]+\.[0-9]+(\.[0-9]+)?$ ]] \
    || die '--version must look like 1.27 or 1.27.0 for archive installation'
  [[ "$ARCHIVE_SHA256" =~ ^[[:xdigit:]]{64}$ ]] \
    || die '--sha256 must be a 64-character hexadecimal checksum for archive installation'
  command -v curl >/dev/null 2>&1 || die 'curl is required for archive installation'
  command -v sha256sum >/dev/null 2>&1 || die 'sha256sum is required for archive verification'
  command -v tar >/dev/null 2>&1 || die 'tar is required for archive installation'
  [[ "$(uname -s)" == 'Linux' ]] || die 'this script supports Linux Go archives only'
  validate_archive_destination

  local arch archive_url archive extracted_go destination_parent
  arch=$(archive_architecture)
  archive_url="https://go.dev/dl/go${GO_VERSION}.linux-${arch}.tar.gz"
  ARCHIVE_TEMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/language-practice-go.XXXXXX")
  archive="$ARCHIVE_TEMP_DIR/go.tar.gz"
  extracted_go="$ARCHIVE_TEMP_DIR/go"
  destination_parent=$(dirname -- "$INSTALL_DIR")

  note "downloading official Go archive: $archive_url"
  curl --fail --location --proto '=https' --tlsv1.2 --output "$archive" "$archive_url"
  note 'verifying Go archive SHA-256 checksum'
  printf '%s  %s\n' "$ARCHIVE_SHA256" "$archive" | sha256sum --check --status - \
    || die 'Go archive checksum verification failed; refusing installation'
  tar --extract --gzip --file "$archive" --directory "$ARCHIVE_TEMP_DIR"
  [[ -x "$extracted_go/bin/go" ]] || die 'verified Go archive did not contain go/bin/go'

  if [[ -e "$INSTALL_DIR" && "$ALLOW_REPLACE" -ne 1 ]]; then
    die "archive destination already exists; pass --replace to replace only: $INSTALL_DIR"
  fi

  if [[ -w "$destination_parent" ]]; then
    note "installing verified Go archive into: $INSTALL_DIR"
    if [[ -e "$INSTALL_DIR" ]]; then rm -rf -- "$INSTALL_DIR"; fi
    mv -- "$extracted_go" "$INSTALL_DIR"
  else
    require_sudo
    note "installing verified Go archive with sudo into: $INSTALL_DIR"
    sudo install -d -- "$destination_parent"
    if [[ -e "$INSTALL_DIR" ]]; then sudo rm -rf -- "$INSTALL_DIR"; fi
    sudo mv -- "$extracted_go" "$INSTALL_DIR"
  fi

  GO_BIN="$INSTALL_DIR/bin/go"
  [[ -x "$GO_BIN" ]] || die "installed Go executable is missing: $GO_BIN"
  GO_INSTALL_STATUS='installed by this run'
}

install_missing_go() {
  if [[ "$NO_INSTALL" -eq 1 || "$INSTALL_METHOD" == 'none' ]]; then
    die 'Go is unavailable and installation was disabled; install Go and rerun this script'
  fi
  detect_install_method
  if [[ "$INSTALL_METHOD" == 'archive' ]]; then
    install_from_archive
  else
    install_from_package
    GO_BIN=$(resolve_command go)
  fi
}

while (($# > 0)); do
  case "$1" in
    --go) (($# >= 2)) || die '--go requires an executable path or command name'; GO_BIN=$2; shift 2 ;;
    --install-method) (($# >= 2)) || die '--install-method requires a value'; INSTALL_METHOD=$2; shift 2 ;;
    --version) (($# >= 2)) || die '--version requires a value'; GO_VERSION=$2; shift 2 ;;
    --sha256) (($# >= 2)) || die '--sha256 requires a value'; ARCHIVE_SHA256=$2; shift 2 ;;
    --install-dir) (($# >= 2)) || die '--install-dir requires a path'; INSTALL_DIR=$2; shift 2 ;;
    --replace) ALLOW_REPLACE=1; shift ;;
    --no-install) NO_INSTALL=1; shift ;;
    --module) (($# >= 2)) || die '--module requires a module path'; [[ -n "$2" ]] || die '--module requires a non-empty module path'; MODULE_PATH=$2; shift 2 ;;
    --tidy) RUN_TIDY=1; shift ;;
    --test) RUN_TEST=1; shift ;;
    --vet) RUN_VET=1; shift ;;
    --help|-h) usage; exit 0 ;;
    *) die "unknown option: $1" ;;
  esac
done

if [[ "$INSTALL_METHOD" == 'archive' && "$NO_INSTALL" -eq 1 ]]; then
  die '--no-install cannot be combined with --install-method archive'
fi
cd -- "$TRACK_DIR"
printf 'Repository: %s\n' "$REPOSITORY_ROOT"
printf 'Go track: %s\n' "$TRACK_DIR"

if [[ -n "$GO_BIN" ]]; then
  GO_BIN=$(resolve_command "$GO_BIN")
elif command -v go >/dev/null 2>&1; then
  GO_BIN=$(command -v go)
else
  install_missing_go
fi

GO_VERSION_OUTPUT=$("$GO_BIN" version)
GO_VERSION_OUTPUT=${GO_VERSION_OUTPUT//$'\n'/ }
printf '%s\n' "$GO_VERSION_OUTPUT"
printf 'Toolchain verification: PASS — executable Go toolchain found at %s.\n' "$GO_BIN"
printf 'Go toolchain installation: %s.\n' "$GO_INSTALL_STATUS"
printf '%s\n' 'Host toolchain version is an observation, not a repository support decision.'
printf '%s\n' 'Go environment (read-only):'
"$GO_BIN" env GOPATH GOMOD GOROOT GOPROXY

if [[ -n "$MODULE_PATH" ]]; then
  if [[ -e go.mod ]]; then
    note 'go.mod already exists; preserving it and skipping go mod init'
  else
    note "initializing project module: go mod init $MODULE_PATH"
    "$GO_BIN" mod init "$MODULE_PATH"
  fi
elif [[ ! -e go.mod ]]; then
  note 'no go.mod found; module initialization is deferred until the learner chooses --module PATH'
fi

if [[ "$RUN_TIDY" -eq 1 ]]; then
  [[ -e go.mod ]] || die '--tidy requires go.mod; pass --module PATH or create the module first'
  note 'running go mod tidy; this may resolve dependencies and access the network'
  "$GO_BIN" mod tidy
fi
if [[ "$RUN_TEST" -eq 1 ]]; then
  note 'running bounded project verification: go test ./...'
  "$GO_BIN" test ./...
fi
if [[ "$RUN_VET" -eq 1 ]]; then
  note 'running bounded project diagnostics: go vet ./...'
  "$GO_BIN" vet ./...
fi

printf '%s\n' 'Bootstrap inspection complete. User-run module, dependency, test, and vet evidence is separate from toolchain verification.'
cat <<'LIMITS'

Cleanup and rollback:
- package-manager changes are system decisions; do not remove them automatically;
- archive cleanup removes only its temporary download/extraction directory;
- archive --replace removes only the explicitly selected --install-dir after the
  user supplied --replace and the archive checksum passed;
- go.mod, go.sum, and project files are learner-owned state and require review
  before manual rollback;
- this script never edits shell profiles or PATH.

Evidence limits:
- static review cannot prove Go installation or archive authenticity beyond the
  checksum command the user actually runs;
- `go version` and `go env` describe this host, not repository support;
- module, test, and vet results prove only the selected local project state;
- no command proves portability, production readiness, race freedom, performance,
  or learner understanding.
LIMITS