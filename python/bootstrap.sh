#!/usr/bin/env bash
# Python bootstrap for this track. Run it from the repository root (or let the
# script locate the repository); it is intentionally not run by the iterator.
#
# Official references:
#   https://docs.python.org/3/library/venv.html
#   https://packaging.python.org/en/latest/guides/installing-using-pip-and-virtual-environments/
#
# The documented primitives are `python -m venv .venv` and `python -m pip`.
# Using the interpreter selected below instead of global `pip` keeps package
# ownership inside the disposable project environment. Activation is only a
# convenience for a user's current shell: `source .venv/bin/activate`, then
# `deactivate`. This script never edits a shell profile.
#
# A venv contains interpreter-specific, disposable state and is not committed
# or treated as movable project state. Recreate it with --recreate when its
# interpreter or layout is stale; recreation removes only the named .venv.
# Dependencies are not invented here: pass an explicit requirements or lock
# file with --requirements, and optionally record the resulting environment with
# --record. A requirements file is an input, not automatically a lockfile.
#
# Local Python 3.14.4 is only a host observation, not a repository support
# decision. The default run installs missing Python host prerequisites using the
# detected Linux package manager; choose and record a version/source policy when
# a distribution package is not appropriate.
set -Eeuo pipefail

usage() {
  cat <<'EOF'
Usage: ./python/bootstrap.sh [options]

Installs missing Python host prerequisites, then creates or reuses a project-
local .venv and performs visible checks. No project dependency is installed
unless --install-packages and --requirements are given.

Options:
  --python PATH          interpreter to use (default: python3, then python)
  --install-method NAME  auto, apt, dnf, pacman, or none (default: auto)
  --no-install           fail instead of installing missing host prerequisites
  --venv PATH            disposable venv path (default: .venv)
  --recreate             remove only the selected venv before creating it
  --requirements PATH    requirements/lock input; not inferred or generated
  --install-packages     run "$PYTHON -m pip install -r PATH" (network may occur)
  --upgrade-pip          explicitly run "$PYTHON -m pip install --upgrade pip"
  --record PATH          explicitly write "$PYTHON -m pip freeze" to PATH
  --help                 show this help

The default installation methods are visible distribution-package commands:
`sudo apt-get install python3 python3-venv python3-pip`,
`sudo dnf install python3 python3-pip`, or
`sudo pacman -S --needed python python-pip`. Use --no-install when the host
toolchain is managed elsewhere. A user-local official build or version manager
such as pyenv remains an explicit alternative. Confirm source, version,
destination, permissions, and rollback limits yourself.
EOF
}

die() { printf 'bootstrap: error: %s\n' "$*" >&2; exit 1; }
note() { printf 'bootstrap: %s\n' "$*"; }

ROOT=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)
cd -- "$ROOT"

PYTHON_BIN=''
VENV='.venv'
REQUIREMENTS=''
RECREATE=0
INSTALL_METHOD=${PYTHON_INSTALL_METHOD:-auto}
INSTALL_TOOLCHAIN=1
INSTALL_PACKAGES=0
UPGRADE_PIP=0
RECORD=''

while (($#)); do
  case $1 in
    --python) (($# >= 2)) || die "--python needs a path"; PYTHON_BIN=$2; shift 2 ;;
    --install-method) (($# >= 2)) || die '--install-method needs auto, apt, dnf, pacman, or none'; INSTALL_METHOD=$2; shift 2 ;;
    --no-install) INSTALL_TOOLCHAIN=0; shift ;;
    --venv) (($# >= 2)) || die "--venv needs a path"; VENV=$2; shift 2 ;;
    --recreate) RECREATE=1; shift ;;
    --requirements) (($# >= 2)) || die "--requirements needs a path"; REQUIREMENTS=$2; shift 2 ;;
    --install-packages) INSTALL_PACKAGES=1; shift ;;
    --upgrade-pip) UPGRADE_PIP=1; shift ;;
    --record) (($# >= 2)) || die "--record needs a path"; RECORD=$2; shift 2 ;;
    --help|-h) usage; exit 0 ;;
    *) die "unknown option: $1 (use --help)" ;;
  esac
done

[[ $VENV != /* ]] || die "--venv must be a path below the repository"
[[ $RECORD != /* ]] || die "--record must be a path below the repository"
VENV=$ROOT/$VENV
[[ $VENV == "$ROOT"/* ]] || die "--venv escaped the repository"
if [[ -n $RECORD ]]; then
  RECORD=$ROOT/$RECORD
  [[ $RECORD == "$ROOT"/* ]] || die "--record escaped the repository"
fi

select_python() {
  if [[ -n $PYTHON_BIN ]]; then return; fi
  if command -v python3 >/dev/null 2>&1; then PYTHON_BIN=$(command -v python3)
  elif command -v python >/dev/null 2>&1; then PYTHON_BIN=$(command -v python)
  fi
}

install_python() {
  ((INSTALL_TOOLCHAIN)) || die 'no Python interpreter found and --no-install was selected'
  case $INSTALL_METHOD in
    auto)
      if command -v apt-get >/dev/null 2>&1; then INSTALL_METHOD=apt
      elif command -v dnf >/dev/null 2>&1; then INSTALL_METHOD=dnf
      elif command -v pacman >/dev/null 2>&1; then INSTALL_METHOD=pacman
      else die 'cannot detect apt-get, dnf, or pacman; install Python or choose --python PATH'; fi ;;
    apt|dnf|pacman|none) ;;
    *) die "unsupported Python install method: $INSTALL_METHOD" ;;
  esac
  [[ $INSTALL_METHOD != none ]] || die 'Python installation disabled; choose --python PATH'
  command -v sudo >/dev/null 2>&1 || die "sudo is required for $INSTALL_METHOD Python installation"
  case $INSTALL_METHOD in
    apt)
      note 'installing Python toolchain: sudo apt-get update'
      sudo apt-get update
      note 'installing Python toolchain: sudo apt-get install -y python3 python3-venv python3-pip'
      sudo apt-get install -y python3 python3-venv python3-pip ;;
    dnf)
      note 'installing Python toolchain: sudo dnf install -y python3 python3-pip'
      sudo dnf install -y python3 python3-pip ;;
    pacman)
      note 'installing Python toolchain: sudo pacman -S --needed --noconfirm python python-pip'
      sudo pacman -S --needed --noconfirm python python-pip ;;
  esac
  PYTHON_BIN=''
  select_python
}

select_python
if [[ -z $PYTHON_BIN ]]; then install_python; fi
[[ -x $PYTHON_BIN ]] || die "interpreter is not executable: $PYTHON_BIN"
if ! "$PYTHON_BIN" -m pip --version >/dev/null 2>&1; then
  install_python
  "$PYTHON_BIN" -m pip --version >/dev/null 2>&1 || die 'Python pip support is unavailable after installation'
fi
if ! "$PYTHON_BIN" -c 'import venv' >/dev/null 2>&1; then
  install_python
  "$PYTHON_BIN" -c 'import venv' >/dev/null 2>&1 || die 'Python venv support is unavailable after installation'
fi
if [[ -n $REQUIREMENTS ]]; then
  [[ $REQUIREMENTS != /* ]] || die "--requirements must be a repository-relative path"
  REQUIREMENTS=$ROOT/$REQUIREMENTS
  [[ $REQUIREMENTS == "$ROOT"/* ]] || die "--requirements escaped the repository"
  [[ -f $REQUIREMENTS ]] || die "requirements input does not exist: $REQUIREMENTS"
fi
if ((INSTALL_PACKAGES || UPGRADE_PIP)) && [[ -z $REQUIREMENTS ]] && ((INSTALL_PACKAGES)); then
  die '--install-packages requires --requirements PATH'
fi

if [[ -e $VENV && ! -d $VENV ]]; then die "venv path is not a directory: $VENV"; fi
if ((RECREATE)) && [[ -d $VENV ]]; then
  case $VENV in "$ROOT/.venv"|"$ROOT"/.venv/*) rm -rf -- "$VENV";;
    *) die 'refusing to recreate a path outside the disposable .venv scope';;
  esac
fi
if [[ ! -d $VENV ]]; then
  note "creating disposable venv with: $PYTHON_BIN -m venv $VENV"
  "$PYTHON_BIN" -m venv "$VENV"
else
  note "reusing existing venv: $VENV (use --recreate only when you intend to discard it)"
fi

VENV_PYTHON=$VENV/bin/python
[[ -x $VENV_PYTHON ]] || die "venv interpreter is missing: $VENV_PYTHON"
note "interpreter: $VENV_PYTHON"
"$VENV_PYTHON" --version
note 'venv import check'
"$VENV_PYTHON" -c 'import sys, venv; assert sys.prefix != sys.base_prefix; print("imports: sys, venv")'
[[ -x $VENV/bin/pip ]] || die "venv pip is missing: $VENV/bin/pip"
note 'pip check (project-local environment)'
"$VENV_PYTHON" -m pip --version
"$VENV_PYTHON" -m pip check

if ((UPGRADE_PIP)); then
  note "network/package operation requested visibly: $VENV_PYTHON -m pip install --upgrade pip"
  "$VENV_PYTHON" -m pip install --upgrade pip
fi
if ((INSTALL_PACKAGES)); then
  note "network/package operation requested visibly: $VENV_PYTHON -m pip install -r $REQUIREMENTS"
  "$VENV_PYTHON" -m pip install -r "$REQUIREMENTS"
fi
if [[ -n $RECORD ]]; then
  mkdir -p -- "$(dirname -- "$RECORD")"
  note "recording dependencies explicitly with: $VENV_PYTHON -m pip freeze"
  "$VENV_PYTHON" -m pip freeze >"$RECORD"
fi

note 'bounded project verification: local pip check completed; no application command was inferred'
note "activate manually in your shell with: source ${VENV#$ROOT}/bin/activate"
note 'deactivate later with: deactivate'
note 'static iterator evidence cannot prove installation, downloads, or project correctness'