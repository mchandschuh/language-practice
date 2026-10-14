#!/usr/bin/env bash
# TypeScript/Node bootstrap for the language-practice track.
#
# This file is a visible, user-run guide. The iterator and reviewer must not run
# it, install software, access npm, or modify the user's host while processing
# this step. Read every command before running it from the repository root.
#
# Source guidance:
# - Node.js installation: https://nodejs.org/en/download/package-manager
# - npm package/lockfile guidance: https://docs.npmjs.com/cli/v11/configuring-npm/package-json
# - TypeScript Handbook: https://www.typescriptlang.org/docs/handbook/
# - TSConfig reference: https://www.typescriptlang.org/tsconfig

set -Eeuo pipefail

readonly TRACK_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
readonly REPOSITORY_ROOT="$(cd -- "$TRACK_DIR/.." && pwd -P)"
cd -- "$TRACK_DIR"

printf 'Repository: %s\n' "$REPOSITORY_ROOT"
printf 'TypeScript track: %s\n' "$TRACK_DIR"
node_version="$(node --version)"
npm_version="$(npm --version)"
printf 'Node: %s\n' "$node_version"
printf 'npm: %s\n' "$npm_version"

cat <<'NOTICE'

Node installation and project dependencies are separate decisions.

Choose and record a Node installation method before continuing. The official
Node package-manager page documents distribution packages and version-manager
options. A distribution package may need sudo and change system state; a
user-local version manager changes the user's toolchain state. This script does
not edit shell profiles or select a repository-supported Node version.

The observed local baseline for this repository was Node 24.21.0 and npm
11.19.0. Those observations are not a support policy. There was no local tsc
when this script was authored.

The commands below intentionally resolve TypeScript from this project. Do not
replace them with a global npm install: a global compiler is an optional,
one-off convenience and is not reproducible project state.
NOTICE

if [[ ! -f package.json ]]; then
  printf '%s\n' 'No package.json found; creating minimal project metadata with npm init -y.'
  npm init -y
else
  printf '%s\n' 'Using existing package.json; it will not be overwritten.'
fi

if [[ ! -f package-lock.json ]]; then
  printf '%s\n' 'No package-lock.json found; installing the local TypeScript dependency will create one.'
else
  printf '%s\n' 'Using existing package-lock.json; inspect it before changing dependencies.'
fi

typescript_install_status='already present'
if ! node -e "const p=require('./package.json'); process.exit(p.devDependencies?.typescript ? 0 : 1)" \
  || [[ ! -x node_modules/.bin/tsc ]]; then
  printf '%s\n' 'Installing TypeScript as a project-local development dependency.'
  printf '%s\n' 'This is the explicit network/package-registry boundary; review the command before allowing it.'
  npm install --save-dev typescript
  typescript_install_status='installed by this run'
else
  printf '%s\n' 'A project-local TypeScript dependency is already declared.'
fi

if [[ ! -f tsconfig.json ]]; then
  cat > tsconfig.json <<'JSON'
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "strict": true,
    "outDir": "dist",
    "rootDir": "src",
    "noEmitOnError": true
  },
  "include": ["src/**/*.ts"]
}
JSON
  printf '%s\n' 'Created a strict NodeNext tsconfig.json example; review it against the project module policy.'
else
  printf '%s\n' 'Using existing tsconfig.json; it will not be overwritten.'
fi

printf '%s\n' 'Local TypeScript:'
typescript_version="$(npx --no-install tsc --version)"
printf '%s\n' "$typescript_version"
printf 'Toolchain verification: PASS — Node.js %s, npm %s, and project-local TypeScript are executable.\n' \
  "$node_version" "$npm_version"
printf 'TypeScript dependency: %s (%s).\n' "$typescript_install_status" "$typescript_version"
printf '%s\n' 'Host toolchain status: Node.js/npm are available; this script reports them but does not install or change them.'
if [[ "$typescript_install_status" == 'installed by this run' ]]; then
  printf '%s\n' 'Toolchain installation: PASS — project-local TypeScript installation completed successfully.'
else
  printf '%s\n' 'Toolchain installation: NOT NEEDED — project-local TypeScript was already installed and verified.'
fi

if [[ ! -d src ]]; then
  mkdir -- src
  printf '%s\n' 'Created empty src/; add learner-owned TypeScript before running the type check.'
fi

printf '%s\n' 'Type checking project sources (noEmit is not assumed; tsconfig controls emission):'
typecheck_output="$(mktemp "${TMPDIR:-/tmp}/language-practice-typescript-typecheck.XXXXXX")"
if npx --no-install tsc --pretty false >"$typecheck_output" 2>&1; then
  cat "$typecheck_output"
  rm -f -- "$typecheck_output"
  printf '%s\n' 'Project type check: PASS.'
else
  typecheck_status=$?
  if grep -q 'error TS18003: No inputs were found' "$typecheck_output" \
    && ! awk '/(^|: )error TS[0-9]+/ && $0 !~ /(^|: )error TS18003: No inputs were found/ { found=1 } END { exit (found ? 0 : 1) }' "$typecheck_output"; then
    rm -f -- "$typecheck_output"
    printf '%s\n' 'Project type check: SKIPPED — no TypeScript source inputs match the current tsconfig include patterns.'
    printf '%s\n' 'Add learner-owned .ts files under the configured source path and rerun to type-check project sources.'
  else
    cat "$typecheck_output" >&2
    rm -f -- "$typecheck_output"
    printf 'Project type check: FAIL (tsc exit %s).\n' "$typecheck_status" >&2
    exit "$typecheck_status"
  fi
fi

runtime_dir="$(mktemp -d "${TMPDIR:-/tmp}/language-practice-typescript.XXXXXX")"
cleanup_runtime() {
  rm -rf -- "$runtime_dir"
}
trap cleanup_runtime EXIT

cat > "$runtime_dir/runtime-check.ts" <<'TS'
export const runtimeCheck = 'TypeScript emitted JavaScript can load in a fresh Node process.';
TS
cat > "$runtime_dir/tsconfig.json" <<'JSON'
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "NodeNext",
    "moduleResolution": "NodeNext",
    "outDir": "out",
    "strict": true
  },
  "files": ["runtime-check.ts"]
}
JSON
npx --no-install tsc --pretty false --project "$runtime_dir/tsconfig.json"
node "$runtime_dir/out/runtime-check.js"
printf '%s\n' 'Fresh-process runtime check passed for the temporary emitted module.'

printf '%s\n' 'Bootstrap verification complete for this host and current project state.'
cat <<'LIMITS'

Evidence limits:
- node/npm version output is host observation, not repository support evidence;
- npx tsc and the checks above prove only this project's current local state;
- a lockfile records dependency resolution but does not prove package authenticity,
  future reproducibility, runtime correctness, or learner understanding;
- the temporary runtime check does not validate the project's own imports or module
  graph.

Cleanup/rollback: the temporary runtime directory is removed automatically. Do not
delete package.json, package-lock.json, tsconfig.json, src/, or node_modules without
reviewing whether they predated this run. Revert only files this bootstrap created
after confirming they are not shared learner work. Global package removal, system
package rollback, and shell-profile changes are manual user decisions.
LIMITS