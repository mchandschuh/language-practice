#!/usr/bin/env bash

set -Eeuo pipefail

repo_root="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd -P)"
test_root=$(mktemp -d "${TMPDIR:-/tmp}/language-practice-temporal-test.XXXXXX")
running_pid=''

cleanup() {
    if [[ -n "$running_pid" ]] && kill -0 "$running_pid" 2>/dev/null; then
        kill -TERM "$running_pid" 2>/dev/null || true
        sleep 1
        kill -KILL "$running_pid" 2>/dev/null || true
    fi
    rm -rf -- "$test_root"
}

trap cleanup EXIT

fail() {
    printf 'FAIL: %s\n' "$*" >&2
    exit 1
}

assert_contains() {
    local actual=$1 expected=$2 message=${3:-}
    [[ "$actual" == *"$expected"* ]] || fail "${message:-expected output to contain $expected}: $actual"
}

assert_file() {
    [[ -e "$1" ]] || fail "expected file to exist: $1"
}

assert_not_file() {
    [[ ! -e "$1" ]] || fail "expected file not to exist: $1"
}

assert_process_gone() {
    local pid=$1
    if kill -0 "$pid" 2>/dev/null; then
        fail "expected process $pid to be stopped"
    fi
}

fake_cli="$test_root/fake-temporal"
cat >"$fake_cli" <<'EOF'
#!/usr/bin/env bash

set -Eeuo pipefail

if [[ "$*" == *'operator cluster health'* ]]; then
    [[ -f "$FAKE_TEMPORAL_HEALTH_FILE" ]]
    exit
fi

if [[ "$1" == server && "$2" == start-dev ]]; then
    : >"$FAKE_TEMPORAL_HEALTH_FILE"
    trap 'rm -f -- "$FAKE_TEMPORAL_HEALTH_FILE"; exit 0' TERM INT
    while :; do
        sleep 1
    done
fi

exit 2
EOF
chmod 0755 "$fake_cli"

fake_systemctl="$test_root/fake-systemctl"
cat >"$fake_systemctl" <<'EOF'
#!/usr/bin/env bash
printf '%s\n' "$*" >>"$FAKE_SYSTEMCTL_LOG"
EOF
chmod 0755 "$fake_systemctl"

export TEMPORAL_CLI="$fake_cli"
export TEMPORAL_SYSTEMCTL="$fake_systemctl"
export FAKE_TEMPORAL_HEALTH_FILE="$test_root/healthy"
export FAKE_SYSTEMCTL_LOG="$test_root/systemctl.log"
export TEMPORAL_DEV_RUNTIME_DIR="$test_root/runtime"
export TEMPORAL_DEV_STATE_DIR="$test_root/state"
export TEMPORAL_DEV_CONFIG_HOME="$test_root/config"
export TEMPORAL_DEV_START_TIMEOUT=3
export TEMPORAL_DEV_STOP_TIMEOUT=3

# shellcheck disable=SC1091
source "$repo_root/init"

if command -v zsh >/dev/null 2>&1; then
    zsh -f -c 'set -u; source "$1"; root=$(_language_practice_repo_root); [[ "$root" == "$2" ]]' \
        zsh "$repo_root/init" "$repo_root"
fi

pid_file="$TEMPORAL_DEV_RUNTIME_DIR/temporal.pid"
log_file="$TEMPORAL_DEV_STATE_DIR/temporal.log"
unit_file="$TEMPORAL_DEV_CONFIG_HOME/systemd/user/language-practice-temporal.service"

if temporal_dev_status >"$test_root/status.out" 2>"$test_root/status.err"; then
    fail 'status unexpectedly succeeded before start'
fi
assert_contains "$(cat "$test_root/status.out")" 'stopped'

start_output=$(temporal_dev_start)
assert_contains "$start_output" 'Started Temporal development server'
assert_file "$pid_file"
assert_file "$log_file"
pid=$(cat "$pid_file")
running_pid=$pid
kill -0 "$pid" 2>/dev/null || fail "started PID is not alive: $pid"

status_output=$(temporal_dev_status)
assert_contains "$status_output" 'running and healthy'
assert_contains "$status_output" "owned PID $pid"

repeat_output=$(temporal_dev_start)
assert_contains "$repeat_output" "already running (owned PID $pid)"
[[ "$(cat "$pid_file")" == "$pid" ]] || fail 'idempotent start changed the owned PID'

temporal_dev_stop
assert_process_gone "$pid"
running_pid=''
assert_not_file "$pid_file"

# A healthy endpoint without this helper's PID file must not be stopped.
: >"$FAKE_TEMPORAL_HEALTH_FILE"
if temporal_dev_start >"$test_root/external.out" 2>"$test_root/external.err"; then
    :
else
    fail 'external healthy endpoint was rejected by start'
fi
assert_contains "$(cat "$test_root/external.out")" 'not owned by this helper'
if temporal_dev_stop >"$test_root/stop-external.out" 2>"$test_root/stop-external.err"; then
    fail 'stop unexpectedly claimed an external endpoint'
fi
assert_contains "$(cat "$test_root/stop-external.err")" 'does not own it'
rm -f -- "$FAKE_TEMPORAL_HEALTH_FILE"

# A dead PID file is repaired without signaling an unrelated process.
stale_pid=999999
while kill -0 "$stale_pid" 2>/dev/null; do
    stale_pid=$((stale_pid + 1))
done
printf '%s\n' "$stale_pid" >"$pid_file"
if temporal_dev_status >"$test_root/stale.out" 2>"$test_root/stale.err"; then
    fail 'stale PID status unexpectedly succeeded'
fi
assert_contains "$(cat "$test_root/stale.out")" 'removed stale PID file'
assert_not_file "$pid_file"

# Missing CLI is actionable and does not create runtime state.
export TEMPORAL_CLI="$test_root/missing-temporal"
if temporal_dev_start >"$test_root/missing.out" 2>"$test_root/missing.err"; then
    fail 'start unexpectedly succeeded without the Temporal CLI'
fi
assert_contains "$(cat "$test_root/missing.err")" 'install_temporal'
export TEMPORAL_CLI="$fake_cli"

# Unit generation is explicit and does not enable or start by default.
temporal_dev_service_install
assert_file "$unit_file"
unit_contents=$(cat "$unit_file")
assert_contains "$unit_contents" 'ExecStart='
assert_contains "$unit_contents" '--port 7233'
assert_contains "$unit_contents" '--ui-port 8233'
assert_contains "$(cat "$FAKE_SYSTEMCTL_LOG")" '--user daemon-reload'
if grep -q -- '--user enable\|--user start' "$FAKE_SYSTEMCTL_LOG"; then
    fail 'default service installation enabled or started the unit'
fi

temporal_dev_service_uninstall
assert_not_file "$unit_file"
assert_contains "$(cat "$FAKE_SYSTEMCTL_LOG")" '--user disable --now language-practice-temporal.service'

temporal_dev_service_install --persistent
assert_file "$unit_file"
[[ -d "$TEMPORAL_DEV_STATE_DIR" ]] || fail 'persistent service did not create its state directory'
unit_contents=$(cat "$unit_file")
assert_contains "$unit_contents" "--db-filename $TEMPORAL_DEV_STATE_DIR/temporal.db"
temporal_dev_service_uninstall
assert_not_file "$unit_file"

printf 'PASS: Temporal init lifecycle and service helper tests\n'