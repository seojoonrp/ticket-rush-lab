#!/usr/bin/env bash
# Re-run every stage N times at its commit and average the results.
#
#   ./bench/run.sh            # all stages, 5 runs each
#   ./bench/run.sh 3-redis    # one stage only
#   RUNS=3 ./bench/run.sh     # change run count
#
# Works by checking out each stage's commit, building+running the server,
# and driving k6. Restores your original branch on exit.
set -euo pipefail

cd "$(dirname "$0")/.."
REPO="$(pwd)"
BENCH="$REPO/bench"
BASE="http://localhost:8080"
RUNS="${RUNS:-5}"           # measured runs per scenario (plus 1 discarded warmup)
RESULTS="$BENCH/results"

# stage : commit : scenarios : drain-queue?(1 only for worker_pool)
STAGES=(
  "1-naive:4667a5b:hotspot,spread:0"
  "2-db_atomic:420cf23:hotspot,spread:0"
  "3-redis:75306a7:hotspot,spread:0"
  "4-worker_pool:c4e1dc9:hotspot,spread:1"
)

LOCK="$BENCH/.lock"
if ! mkdir "$LOCK" 2>/dev/null; then
  echo "another ./bench/run.sh is already running (lock: $LOCK)." >&2
  echo "if you're sure none is, remove it:  rmdir $LOCK" >&2
  exit 1
fi

ORIG="$(git rev-parse --abbrev-ref HEAD)"; [ "$ORIG" = HEAD ] && ORIG="$(git rev-parse HEAD)"
SERVER_PID=""
cleanup() {
  [ -n "$SERVER_PID" ] && kill "$SERVER_PID" 2>/dev/null || true
  git checkout -q "$ORIG" 2>/dev/null || true
  rmdir "$LOCK" 2>/dev/null || true
}
trap cleanup EXIT

# --- infra (mongo + redis) ---
echo ">> starting mongo + redis"
docker compose up -d
for _ in $(seq 60); do
  if docker exec ticket-rush-mongo mongosh --quiet --eval 'rs.status().myState' 2>/dev/null | grep -q '^1$'; then break; fi
  sleep 1
done

start_server() {
  go build -o "$BENCH/server" ./cmd/server
  # env var names differ across stages (naive/db_atomic use DB_URI/DB_NAME;
  # redis+ use MONGO_URI/MONGO_DB_NAME/REDIS_ADDR/WORKER_COUNT/BUFFER_SIZE).
  # export a superset so every commit finds what it reads (godotenv won't override these).
  DB_URI="mongodb://localhost:27017" DB_NAME="ticket-rush-lab" \
  MONGO_URI="mongodb://localhost:27017" MONGO_DB_NAME="ticket-rush-lab" \
  REDIS_ADDR="localhost:6379" WORKER_COUNT=32 BUFFER_SIZE=256 \
    "$BENCH/server" >"$BENCH/server.log" 2>&1 &
  SERVER_PID=$!
  for _ in $(seq 60); do
    if ! kill -0 "$SERVER_PID" 2>/dev/null; then
      echo "!! server crashed on startup:" >&2; tail -3 "$BENCH/server.log" >&2; exit 1
    fi
    if curl -s --max-time 2 -o /dev/null "$BASE/shows"; then return 0; fi
    sleep 0.3
  done
  echo "server never came up; see bench/server.log" >&2; exit 1
}
stop_server() {
  [ -n "$SERVER_PID" ] && kill "$SERVER_PID" 2>/dev/null || true
  wait "$SERVER_PID" 2>/dev/null || true
  SERVER_PID=""; sleep 1
}

one_run() {   # stage scenario idx drain
  local stage="$1" scn="$2" idx="$3" drain="$4"
  local dir="$RESULTS/$stage/$scn"; mkdir -p "$dir"
  local sj="$dir/run-$idx.summary.json" vj="$dir/run-$idx.verify.json"
  local show
  if [ "$scn" = hotspot ]; then
    local c; c="$(curl -s --max-time 10 -X POST "$BASE/shows" -H 'Content-Type: application/json' -d '{"seatCount":1}')"
    show="$(jq -r .show.id <<<"$c")"; local seat; seat="$(jq -r '.seats[0].id' <<<"$c")"
    k6 run -e BASE="$BASE" -e SHOW_ID="$show" -e SEAT_ID="$seat" -e ATTACKERS=500 \
           -e SUMMARY_OUT="$sj" "$BENCH/wrap_hotspot.js"
  else
    local c; c="$(curl -s --max-time 10 -X POST "$BASE/shows" -H 'Content-Type: application/json' -d '{"seatCount":100}')"
    show="$(jq -r .show.id <<<"$c")"; local seats; seats="$(jq -r '[.seats[].id]|join(",")' <<<"$c")"
    k6 run -e BASE="$BASE" -e SEAT_IDS="$seats" -e PEAK_VUS=500 \
           -e SUMMARY_OUT="$sj" "$BENCH/wrap_spread.js"
  fi
  if [ "$drain" = 1 ]; then
    for _ in $(seq 300); do
      if [ "$(curl -s --max-time 2 "$BASE/health/queue" | jq .depth 2>/dev/null)" = 0 ]; then break; fi
      sleep 0.2
    done
  fi
  curl -s --max-time 10 "$BASE/shows/$show/verify" >"$vj"
}

rm -rf "$RESULTS"
for e in "${STAGES[@]}"; do
  IFS=: read -r stage commit scns drain <<<"$e"
  if [ -n "${1:-}" ] && [ "$1" != "$stage" ]; then continue; fi
  echo ">> === $stage ($commit) ==="
  git checkout -q "$commit"
  start_server
  IFS=, read -ra arr <<<"$scns"
  for scn in "${arr[@]}"; do
    echo ">> $stage/$scn: warmup + $RUNS runs"
    one_run "$stage" "$scn" 0 "$drain"                    # warmup, ignored by aggregate
    for i in $(seq 1 "$RUNS"); do one_run "$stage" "$scn" "$i" "$drain"; done
  done
  stop_server
done

git checkout -q "$ORIG"
echo ">> aggregating"
python3 "$BENCH/aggregate.py" "$RESULTS" "$BENCH/out"
