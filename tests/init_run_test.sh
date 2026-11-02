#!/usr/bin/env bash

set -Eeuo pipefail

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)"
test_root=$(mktemp -d "${TMPDIR:-/tmp}/language-practice-run-test.XXXXXX")
running_pid=''

cleanup() {
    if [[ -n "$running_pid" ]] && kill -0 -- -"$running_pid" 2>/dev/null; then
        kill -TERM -- -"$running_pid" 2>/dev/null || true
        sleep 1
        kill -KILL -- -"$running_pid" 2>/dev/null || true
    fi
    rm -rf -- "$test_root"
}

trap cleanup EXIT

fail() {
    printf 'FAIL: %s\n' "$*" >&2
    exit 1
}

assert_contains() {
    local actual=$1 expected=$2
    [[ "$actual" == *"$expected"* ]] || fail "expected output to contain $expected: $actual"
}

fake_go="$test_root/go"
cat >"$fake_go" <<'EOF'
#!/usr/bin/env bash

set -Eeuo pipefail

printf '%s\n' "$@" >"$RUN_FAKE_ARGS_FILE"

if [[ "${RUN_FAKE_MODE:-}" == background ]]; then
    bash -c 'trap "exit 0" TERM INT; while :; do read -r -t 1 _ || true; done' &
    child_pid=$!
    : >"$RUN_FAKE_STARTED_FILE"
    trap 'kill -TERM "$child_pid" 2>/dev/null || true; exit 0' TERM INT
    while :; do
        read -r -t 1 _ || true
    done
fi

printf 'foreground-ok\n'
EOF
chmod 0755 "$fake_go"

export PATH="$test_root:$PATH"
export RUN_FAKE_ARGS_FILE="$test_root/args"
export RUN_FAKE_STARTED_FILE="$test_root/started"

# shellcheck disable=SC1091
source "$repo_root/init"

foreground_output=$(run positiveNumbers --num 7)
assert_contains "$foreground_output" 'foreground-ok'
mapfile -t foreground_args <"$RUN_FAKE_ARGS_FILE"
[[ "${foreground_args[*]}" == 'run ./basic/positiveNumbers --num 7' ]] \
    || fail "foreground args = ${foreground_args[*]}"

if run --unknown positiveNumbers >/dev/null 2>"$test_root/invalid.err"; then
    fail 'unknown run option unexpectedly succeeded'
fi
assert_contains "$(cat "$test_root/invalid.err")" 'unknown run option: --unknown'

if run --background missing-app >"$test_root/missing.out" 2>&1; then
    fail 'background run with an unknown app unexpectedly succeeded'
fi
assert_contains "$(cat "$test_root/missing.out")" 'unknown app: missing-app'

run_output_file="$test_root/run.out"
RUN_FAKE_MODE=background run -b positiveNumbers --num 9 >"$run_output_file"
run_output=$(cat "$run_output_file")
assert_contains "$run_output" 'Started positiveNumbers in background process group (PID '
pid=$(sed -n 's/.*PID \([0-9][0-9]*\)).*/\1/p' "$run_output_file")
[[ -n "$pid" ]] || fail "could not extract background PID: $run_output"
running_pid=$pid

for _ in {1..30}; do
    [[ -e "$RUN_FAKE_STARTED_FILE" ]] && break
    sleep 0.1
done
assert_contains "$(cat "$RUN_FAKE_ARGS_FILE")" './basic/positiveNumbers'
grep -Fx -- '--num' "$RUN_FAKE_ARGS_FILE" >/dev/null || fail 'background args omitted --num'
grep -Fx -- '9' "$RUN_FAKE_ARGS_FILE" >/dev/null || fail 'background args omitted value'
kill -0 -- -"$pid" 2>/dev/null || fail "background process group is not alive: $pid"

if zsh -f -c 'source "$1"; [[ "$(run --help 2>&1)" == *"unknown run option"* ]]' zsh "$repo_root/init"; then
    :
else
    fail 'zsh did not preserve run option parsing'
fi

kill -TERM -- -"$pid"
for _ in {1..30}; do
    if ! kill -0 -- -"$pid" 2>/dev/null; then
        running_pid=''
        break
    fi
    sleep 0.1
done
[[ -z "$running_pid" ]] || fail "background process did not stop: $pid"

printf 'PASS: init run foreground/background tests\n'