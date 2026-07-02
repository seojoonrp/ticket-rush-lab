#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE="${BASE:-http://localhost:8080}"
ATTACKERS="${ATTACKERS:-500}" # concurrent request count
SEAT_COUNT="${SEAT_COUNT:-1}"

echo "1) Register show"
CREATE_RES=$(curl -s -X POST "${BASE}/shows" \
  -H 'Content-Type: application/json' \
  -d "{\"seatCount\": ${SEAT_COUNT}}")

SHOW_ID=$(echo "${CREATE_RES}" | jq -r '.show.id')
SEAT_ID=$(echo "${CREATE_RES}" | jq -r '.seats[0].id')

if [[ -z "${SEAT_ID}" || "${SEAT_ID}" == "null" ]]; then
  echo "[ERROR] failed to extract SEAT_ID" >&2
  exit 1
fi

echo "  show_id = ${SHOW_ID}"
echo "  seat_id = ${SEAT_ID}"
echo

echo "2) Execute ${ATTACKERS} concurrent requests"
k6 run \
  -e BASE="${BASE}" \
  -e SHOW_ID="${SHOW_ID}" \
  -e SEAT_ID="${SEAT_ID}" \
  -e ATTACKERS="${ATTACKERS}" \
  "${SCRIPT_DIR}/hotspot.js"
echo

echo "Draining write queue..."
until [ "$(curl -s "${BASE}/health/queue" | jq '.depth')" -eq 0 ]; do
  sleep 0.2
done
echo

echo "3) Verification"
VERIFY_RES=$(curl -s "${BASE}/shows/${SHOW_ID}/verify")
echo "${VERIFY_RES}" | jq .
echo