import http from "k6/http";
import { check } from "k6";
import { Counter } from "k6/metrics";

const booked = new Counter("booked_success");

const BASE = __ENV.BASE || "http://localhost:8080";
const SEAT_ID = __ENV.SEAT_ID;
const SHOW_ID = __ENV.SHOW_ID;
const ATTACKERS = Number(__ENV.ATTACKERS || 500); // 동시 요청 수

export const options = {
  scenarios: {
    hotspot: {
      executor: "shared-iterations",
      vus: ATTACKERS,
      iterations: ATTACKERS,
      maxDuration: "30s",
    },
  },
};

export default function () {
  if (!SEAT_ID) {
    throw new Error("SEAT_ID is needed");
  }

  const userID = `user_${__VU}`;

  const res = http.post(`${BASE}/seats/${SEAT_ID}/book`, null, {
    headers: { "X-User-ID": userID },
  });

  check(res, {
    "status is 200 or 409": (r) => r.status === 200 || r.status === 409,
  });

  if (res.status === 200) {
    booked.add(1);
  }
}
