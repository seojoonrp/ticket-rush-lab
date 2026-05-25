import http from "k6/http";
import { check } from "k6";
import { Counter } from "k6/metrics";

const booked = new Counter("booked_success");

const BASE = __ENV.BASE || "http://localhost:8080";
const SEAT_IDS = (__ENV.SEAT_IDS || "").split(",").filter(Boolean);

export const options = {
  scenarios: {
    spread: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [
        { duration: "10s", target: Number(__ENV.PEAK_VUS || 500) },
        { duration: "20s", target: Number(__ENV.PEAK_VUS || 500) },
        { duration: "5s", target: 0 },
      ],
    },
  },
  thresholds: {
    // p95 latency
    http_req_duration: ["p(95)<500"],
  },
};

export default function () {
  if (SEAT_IDS.length === 0) {
    throw new Error("SEAT_ID is needed");
  }

  const seatID = SEAT_IDS[Math.floor(Math.random() * SEAT_IDS.length)];
  const userID = `user_${__VU}_${__ITER}`;

  const res = http.post(`${BASE}/seats/${seatID}/book`, null, {
    headers: { "X-User-ID": userID },
  });

  check(res, {
    "status is 200 or 409": (r) => r.status === 200 || r.status === 409,
  });

  if (res.status === 200) booked.add(1);
}
