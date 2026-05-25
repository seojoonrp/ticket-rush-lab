#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE="${BASE:-http://localhost:8080}"
SEAT_COUNT="${SEAT_COUNT:-100}" # seat count
PEAK_VUS="${PEAK_VUS:-500}"

echo "1) Register show with ${SEAT_COUNT} seats"
CREATE_RES=$(curl -s -X POST "${BASE}/shows" \
  -H 'Content-Type: application/json' \
  -d "{\"seatCount\": ${SEAT_COUNT}}")

SHOW_ID=$(echo "${CREATE_RES}" | jq -r '.show.id')
SEAT_IDS=$(echo "${CREATE_RES}" | jq -r '[.seats[].id] | join(",")')

if [[ -z "${SEAT_IDS}" || "${SEAT_IDS}" == "null" ]]; then
  echo "[ERROR] failed to extract SEAT_ID" >&2
  exit 1
fi

echo "  show_id = ${SHOW_ID}"
echo

echo "2) Execute requests on ${SEAT_COUNT} different seats (max ${PEAK_VUS} VU)"
k6 run \
  -e BASE="${BASE}" \
  -e SEAT_IDS="${SEAT_IDS}" \
  -e PEAK_VUS="${PEAK_VUS}" \
  "${SCRIPT_DIR}/spread.js"
echo

echo "3) Verification"
VERIFY_RES=$(curl -s "${BASE}/shows/${SHOW_ID}/verify")
echo "${VERIFY_RES}" | jq .
echo