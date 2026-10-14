#!/usr/bin/env bash
# C++/CMake bootstrap for the language-practice track.
#
# This is a visible, user-run guide. The iterator and reviewer must only inspect
# and statically validate it. They must not run this script, install packages,
# invoke sudo, run CMake or a compiler, create a build directory, access the
# network, or perform cleanup while processing WI-0004.
#
# Official references:
# - GCC: https://gcc.gnu.org/onlinedocs/
# - Clang: https://clang.llvm.org/docs/GettingStarted.html
# - CMake tutorial: https://cmake.org/cmake/help/latest/guide/tutorial/index.html
# - CMake presets: https://cmake.org/cmake/help/latest/manual/cmake-presets.7.html
# - ISO C++ getting started: https://isocpp.org/get-started
#
# The repository has no selected C++ standard, compiler, CMake generator, or
# dependency policy. Local GCC 15.2 and Clang 21.1.8 are observations only. The
# default C++17 choice below is a bootstrap default, not repository support.
# Choose and record the compiler and minimum standard before exercises rely on
# concepts, ranges, jthread, coroutines, or other standard-version features.
# Compiler verification is equivalent to `g++ --version` or `clang++ --version`.
#
# Automated installation is limited to visible distribution package commands.
# Compiler source builds, vendor archives, and user-local LLVM/GCC installations
# are intentionally not automated. Install those separately, put the chosen
# tools on the current session PATH, and rerun with --no-install if needed.

set -Eeuo pipefail

readonly TRACK_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
readonly REPOSITORY_ROOT="$(cd -- "$TRACK_DIR/.." && pwd -P)"
readonly INVOCATION_DIR="$(pwd -P)"

COMPILER_REQUEST="${CXX:-g++}"
COMPILER_PATH=''
CXX_STANDARD='17'
INSTALL_METHOD='auto'
INSTALL_ENABLED=1
CHECK_PROJECT=0
SOURCE_DIR="$TRACK_DIR"
BUILD_DIR="$TRACK_DIR/build"
BUILD_JOBS=1

die() {
  printf 'bootstrap error: %s\n' "$*" >&2
  exit 1
}

note() {
  printf '[cpp/bootstrap] %s\n' "$*"
}

usage() {
  cat <<'EOF'
Usage: ./cpp/bootstrap.sh [options]

Reports or installs a C++ compiler and CMake. Project configure/build/test is
opt-in and always uses an out-of-source build directory.

Toolchain options:
  --compiler NAME|PATH    g++ or clang++ (default: CXX or g++)
  --standard N            C++ standard 17, 20, or 23 (default: 17)
  --install-method NAME   auto, apt, dnf, pacman, or none (default: auto)
  --no-install            fail instead of installing missing prerequisites

Project verification options:
  --check-project         configure, build, and run CTest for the project
  --source-dir PATH       CMake source directory (default: ./cpp)
  --build-dir PATH        disposable out-of-source build directory
                          (default: ./cpp/build)
  --build-jobs N          bounded build parallelism (default: 1)
  --help                  show this help

Without --check-project, this script only reports the selected toolchain. It
does not create a build directory or modify project files.

CMakePresets.json may be committed when its compiler, standard, generator, and
build directory are repository decisions. CMakeUserPresets.json is for local
developer settings and should not be used as a shared support contract. A
toolchain file is an explicit CMake input; this script does not create or edit
one.
EOF
}

resolve_executable() {
  local requested=$1
  if [[ "$requested" == */* ]]; then
    [[ -x "$requested" ]] || die "compiler is not executable: $requested"
    printf '%s\n' "$requested"
  else
    command -v -- "$requested" || die "compiler command not found: $requested"
  fi
}

compiler_family() {
  case "$(basename -- "$COMPILER_REQUEST")" in
    clang++|clang*) printf '%s\n' 'clang' ;;
    g++|gcc*) printf '%s\n' 'gcc' ;;
    *) die "unsupported compiler selection: $COMPILER_REQUEST; choose g++ or clang++" ;;
  esac
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
        die 'cannot detect apt-get, dnf, or pacman; install prerequisites manually or choose --install-method none'
      fi
      ;;
    apt|dnf|pacman|none) ;;
    *) die "unsupported --install-method: $INSTALL_METHOD" ;;
  esac
}

require_sudo() {
  command -v sudo >/dev/null 2>&1 || die "sudo is required for $INSTALL_METHOD installation"
}

install_prerequisites() {
  ((INSTALL_ENABLED)) || die 'selected compiler or cmake is missing and --no-install was selected'
  detect_install_method
  [[ "$INSTALL_METHOD" != 'none' ]] || die 'compiler or cmake is missing and installation method is none'
  require_sudo

  local family
  family=$(compiler_family)
  case "$INSTALL_METHOD:$family" in
    apt:gcc)
      note 'installing C++ host prerequisites: sudo apt-get update'
      sudo apt-get update
      note 'installing GCC and CMake: sudo apt-get install -y build-essential cmake'
      sudo apt-get install -y build-essential cmake
      ;;
    apt:clang)
      note 'installing Clang and CMake: sudo apt-get update'
      sudo apt-get update
      note 'installing Clang and CMake: sudo apt-get install -y clang cmake'
      sudo apt-get install -y clang cmake
      ;;
    dnf:gcc)
      note 'installing GCC and CMake: sudo dnf install -y gcc-c++ cmake'
      sudo dnf install -y gcc-c++ cmake
      ;;
    dnf:clang)
      note 'installing Clang and CMake: sudo dnf install -y clang cmake'
      sudo dnf install -y clang cmake
      ;;
    pacman:*)
      note 'installing C++ tools: sudo pacman -S --needed --noconfirm base-devel cmake'
      sudo pacman -S --needed --noconfirm base-devel cmake
      ;;
    *) die "no package mapping for $INSTALL_METHOD and $family" ;;
  esac
}

validate_options() {
  [[ "$CXX_STANDARD" =~ ^(17|20|23)$ ]] \
    || die '--standard must be 17, 20, or 23'
  [[ "$BUILD_JOBS" =~ ^[1-9][0-9]*$ ]] \
    || die '--build-jobs must be a positive integer'
  [[ -d "$SOURCE_DIR" ]] || die "source directory does not exist: $SOURCE_DIR"
  if [[ "$CHECK_PROJECT" -eq 1 ]]; then
    [[ -f "$SOURCE_DIR/CMakeLists.txt" ]] \
      || die "--check-project requires CMakeLists.txt in source directory: $SOURCE_DIR"
  fi
}

verify_toolchain() {
  COMPILER_PATH=$(resolve_executable "$COMPILER_REQUEST")
  command -v cmake >/dev/null 2>&1 \
    || die 'cmake is still unavailable after prerequisite installation'
  printf '%s\n' 'Compiler version:'
  "$COMPILER_PATH" --version
  printf '%s\n' 'CMake version:'
  cmake --version
  printf 'Toolchain verification: PASS — %s and CMake are executable.\n' "$COMPILER_PATH"
  printf 'Selected C++ standard: C++%s (bootstrap default/choice, not repository support).\n' "$CXX_STANDARD"
}

check_project() {
  [[ "$CHECK_PROJECT" -eq 1 ]] || return 0
  note "configuring out-of-source build directory: $BUILD_DIR"
  cmake -S "$SOURCE_DIR" -B "$BUILD_DIR" \
    -DCMAKE_CXX_COMPILER="$COMPILER_PATH" \
    -DCMAKE_CXX_STANDARD="$CXX_STANDARD" \
    -DCMAKE_CXX_STANDARD_REQUIRED=ON \
    -DCMAKE_CXX_EXTENSIONS=OFF

  note 'verifying configured C++ standard'
  local configured_standard
  configured_standard=$(cmake -LA -N -B "$BUILD_DIR" \
    | awk -F= '$1 ~ /^CMAKE_CXX_STANDARD:/ { print $2; exit }')
  [[ "$configured_standard" == "$CXX_STANDARD" ]] \
    || die "configured C++ standard was '$configured_standard', expected '$CXX_STANDARD'"

  note "building with bounded parallelism: $BUILD_JOBS job(s)"
  cmake --build "$BUILD_DIR" --parallel "$BUILD_JOBS"
  note 'running bounded CTest verification'
  ctest --test-dir "$BUILD_DIR" --output-on-failure
  printf '%s\n' 'Project configure/build/test verification: PASS.'
}

normalize_path() {
  local requested=$1
  if [[ "$requested" == /* ]]; then
    printf '%s\n' "$requested"
  else
    printf '%s/%s\n' "$INVOCATION_DIR" "$requested"
  fi
}

while (($# > 0)); do
  case "$1" in
    --compiler) (($# >= 2)) || die '--compiler requires g++ or clang++'; COMPILER_REQUEST=$2; shift 2 ;;
    --standard) (($# >= 2)) || die '--standard requires 17, 20, or 23'; CXX_STANDARD=$2; shift 2 ;;
    --install-method) (($# >= 2)) || die '--install-method requires auto, apt, dnf, pacman, or none'; INSTALL_METHOD=$2; shift 2 ;;
    --no-install) INSTALL_ENABLED=0; shift ;;
    --check-project) CHECK_PROJECT=1; shift ;;
    --source-dir) (($# >= 2)) || die '--source-dir requires a path'; SOURCE_DIR=$2; shift 2 ;;
    --build-dir) (($# >= 2)) || die '--build-dir requires a path'; BUILD_DIR=$2; shift 2 ;;
    --build-jobs) (($# >= 2)) || die '--build-jobs requires a positive integer'; BUILD_JOBS=$2; shift 2 ;;
    --help|-h) usage; exit 0 ;;
    *) die "unknown option: $1" ;;
  esac
done

if [[ "$COMPILER_REQUEST" == */* ]]; then
  COMPILER_REQUEST=$(normalize_path "$COMPILER_REQUEST")
fi
SOURCE_DIR=$(normalize_path "$SOURCE_DIR")
BUILD_DIR=$(normalize_path "$BUILD_DIR")
cd -- "$TRACK_DIR"
printf 'Repository: %s\n' "$REPOSITORY_ROOT"
printf 'C++ track: %s\n' "$TRACK_DIR"
validate_options

if ! command -v -- "$COMPILER_REQUEST" >/dev/null 2>&1 || ! command -v cmake >/dev/null 2>&1; then
  install_prerequisites
fi

verify_toolchain
check_project

printf '%s\n' 'Bootstrap inspection complete. Toolchain, project, portability, and learning evidence remain separate.'
cat <<'LIMITS'

Cleanup and rollback:
- package-manager changes are system decisions; do not remove them automatically;
- build output is confined to the explicitly selected out-of-source build directory;
- remove only that named build directory manually after confirming it contains no
  learner-owned files;
- CMakePresets.json and project toolchain files are not created or edited here;
- this script never edits shell profiles or PATH.

Evidence limits:
- compiler/CMake version output describes this host, not repository support;
- package-manager success does not prove compiler features, project correctness,
  portability, performance, or learner understanding;
- configure/build/CTest results, when the user opts in, prove only that project
  and host snapshot;
- no static check proves another Linux host can reproduce the setup.
LIMITS