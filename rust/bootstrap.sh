#!/usr/bin/env bash
# Rust bootstrap for the language-practice track.
#
# Official references:
#   https://www.rust-lang.org/tools/install
#   https://rust-lang.github.io/rustup/
#   https://doc.rust-lang.org/cargo/
#
# This script is user-run documentation, not an iterator action. It does not
# edit a shell profile. The repository has not selected a supported Rust
# channel or edition; stable is the explicit, changeable bootstrap default,
# not a repository support claim. Rust Book 1.97 / Edition 2024 references are
# context for the learning plan, not an execution contract.
set -Eeuo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)
REPO_ROOT=$(cd -- "$SCRIPT_DIR/.." && pwd -P)
RUST_ROOT=$SCRIPT_DIR
TOOLCHAIN=${RUST_TOOLCHAIN:-stable}
INSTALL_METHOD=auto
INSTALL_TOOLCHAIN=1
MAKE_DEFAULT=0
PROJECT_PATH=
RUN_CHECK=0
RUN_TEST=0
RUN_FMT=0
RUN_CLIPPY=0
COMPONENTS=()
TARGETS=()

die() { printf 'error: %s\n' "$*" >&2; exit 1; }
usage() {
  cat <<'EOF'
Usage: ./rust/bootstrap.sh [options]

Host and toolchain setup:
  --toolchain CHANNEL       rustup channel (default: stable, not support policy)
  --install-method METHOD   auto, apt, dnf, pacman, or none (default: auto)
  --install-toolchain       install CHANNEL with rustup (default)
  --no-install-toolchain    inspect an existing rustup installation only
  --default-toolchain       explicitly run `rustup default CHANNEL`
  --component NAME          install a rustup component for CHANNEL
  --target TRIPLE           install a rustup target for CHANNEL

Project and checks:
  --new-project PATH         run `cargo new PATH`; PATH must be below ./rust
  --check                    run bounded `cargo check --locked`
  --test                     run bounded `cargo test --locked`
  --fmt                      run bounded `cargo fmt --check`
  --clippy                   run bounded `cargo clippy --locked`
  --help                    show this help

The default path installs missing curl/sha256sum prerequisites through a
detected distribution package manager and explicit sudo. It then installs
rustup without changing the user's default toolchain. Use
--default-toolchain only after inspecting the consequences of that global
selection. Package installation, downloads, and Cargo commands are visible
user actions and are never run by the repository iterator.
EOF
}

while (($#)); do
  case $1 in
    --toolchain) (($# >= 2)) || die '--toolchain needs a channel'; TOOLCHAIN=$2; shift 2 ;;
    --install-method) (($# >= 2)) || die '--install-method needs auto, apt, dnf, pacman, or none'; INSTALL_METHOD=$2; shift 2 ;;
    --install-toolchain) INSTALL_TOOLCHAIN=1; shift ;;
    --no-install-toolchain) INSTALL_TOOLCHAIN=0; shift ;;
    --default-toolchain) MAKE_DEFAULT=1; shift ;;
    --component) (($# >= 2)) || die '--component needs a name'; COMPONENTS+=("$2"); shift 2 ;;
    --target) (($# >= 2)) || die '--target needs a triple'; TARGETS+=("$2"); shift 2 ;;
    --new-project) (($# >= 2)) || die '--new-project needs a path'; PROJECT_PATH=$2; shift 2 ;;
    --check) RUN_CHECK=1; shift ;;
    --test) RUN_TEST=1; shift ;;
    --fmt) RUN_FMT=1; shift ;;
    --clippy) RUN_CLIPPY=1; shift ;;
    --help|-h) usage; exit 0 ;;
    *) die "unknown option: $1" ;;
  esac
done

case $INSTALL_METHOD in auto|apt|dnf|pacman|none) ;; *) die "unsupported --install-method: $INSTALL_METHOD" ;; esac
[[ $TOOLCHAIN != -* && $TOOLCHAIN != *[[:space:]]* ]] || die 'toolchain channel must be a non-option, single token'

have() { command -v "$1" >/dev/null 2>&1; }
install_prerequisites() {
  local manager=$INSTALL_METHOD
  if [[ $manager == auto ]]; then
    if have apt-get; then manager=apt
    elif have dnf; then manager=dnf
    elif have pacman; then manager=pacman
    else die 'cannot detect apt-get, dnf, or pacman; install curl and sha256sum manually or choose --install-method none'; fi
  fi
  [[ $manager != none ]] || die 'prerequisite installation disabled; install curl and sha256sum manually'
  have sudo || die "sudo is required for explicit $manager prerequisite installation"
  case $manager in
    apt) sudo apt-get update && sudo apt-get install -y curl coreutils ;;
    dnf) sudo dnf install -y curl coreutils ;;
    pacman) sudo pacman -Sy --needed --noconfirm curl coreutils ;;
  esac
}

if ! have curl || ! have sha256sum; then
  install_prerequisites
fi

have rustup || {
  tmpdir=$(mktemp -d "${TMPDIR:-/tmp}/rust-bootstrap.XXXXXX")
  trap 'rm -rf -- "$tmpdir"' EXIT
  installer=$tmpdir/rustup-init
  printf 'Downloading the official rustup installer to %s\n' "$installer"
  curl --fail --location --proto '=https' --tlsv1.2 --silent --show-error \
    https://sh.rustup.rs --output "$installer"
  chmod 700 "$installer"
  printf 'Installer SHA-256 (record this with your host evidence): '
  sha256sum "$installer"
  "$installer" -y --no-modify-path --default-toolchain none
  export PATH="$HOME/.cargo/bin:$PATH"
}
have rustup || die 'rustup is still unavailable after installation'

printf '%s\n' 'Rustup host/toolchain inspection:'
rustup show
rustup toolchain list
rustup component list --installed --toolchain "$TOOLCHAIN" || true
rustup target list --installed --toolchain "$TOOLCHAIN" || true

if ((INSTALL_TOOLCHAIN)); then
  rustup toolchain install "$TOOLCHAIN"
fi
if ((MAKE_DEFAULT)); then
  printf 'Changing the user default toolchain to %s (explicit --default-toolchain request).\n' "$TOOLCHAIN"
  rustup default "$TOOLCHAIN"
fi
for component in "${COMPONENTS[@]}"; do
  rustup component add "$component" --toolchain "$TOOLCHAIN"
done
for target in "${TARGETS[@]}"; do
  rustup target add "$target" --toolchain "$TOOLCHAIN"
done

if [[ -n $PROJECT_PATH ]]; then
  case $PROJECT_PATH in /*) candidate=$PROJECT_PATH ;; *) candidate=$REPO_ROOT/$PROJECT_PATH ;; esac
  candidate=$(realpath -m -- "$candidate")
  rust_root_real=$(realpath -m -- "$RUST_ROOT")
  [[ $candidate == "$rust_root_real"/* ]] || die '--new-project must remain below ./rust'
  [[ ! -e $candidate ]] || die "project path already exists: $candidate"
  cargo "+$TOOLCHAIN" new "$candidate"
fi

if ((RUN_CHECK || RUN_TEST || RUN_FMT || RUN_CLIPPY)); then
  [[ -n $PROJECT_PATH ]] || die 'project checks require --new-project PATH or a future project path option'
  case $PROJECT_PATH in /*) project=$PROJECT_PATH ;; *) project=$REPO_ROOT/$PROJECT_PATH ;; esac
  project=$(realpath -m -- "$project")
  [[ -f $project/Cargo.toml ]] || die "Cargo.toml not found: $project"
  [[ -f $project/Cargo.lock ]] || die "Cargo.lock not found; create/review it before locked checks: $project"
  (cd -- "$project" && {
    ((RUN_CHECK)) && cargo "+$TOOLCHAIN" check --locked
    ((RUN_TEST)) && cargo "+$TOOLCHAIN" test --locked
    ((RUN_FMT)) && cargo "+$TOOLCHAIN" fmt --check
    ((RUN_CLIPPY)) && cargo "+$TOOLCHAIN" clippy --locked
  })
fi

cat <<'EOF'
Bootstrap complete for the requested visible stages. Cargo target/ caches and
the rustup installation are user/project state; inspect them before cleanup.
This script does not remove unrelated state, Cargo.lock, toolchains, or shell
profiles. Pin a toolchain (for example with rust-toolchain.toml) and choose an
edition deliberately when the repository makes that decision; neither choice
is currently declared by this repository.
EOF