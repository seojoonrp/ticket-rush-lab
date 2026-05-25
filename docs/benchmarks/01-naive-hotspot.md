```bash
1) Register show
  show_id = 6a14654ba90b3b85d18e30d7
  seat_id = 6a14654ba90b3b85d18e30d8

2) Execute 500 concurrent requests

         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/


     execution: local
        script: /home/seojoonrp/practice/api_servers/ticket-rush-lab/loadtest/hotspot.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m0s max duration (incl. graceful stop):
              * hotspot: 500 iterations shared among 500 VUs (maxDuration: 30s, gracefulStop: 30s)



  █ TOTAL RESULTS

    checks_total.......: 500     4444.303688/s
    checks_succeeded...: 100.00% 500 out of 500
    checks_failed......: 0.00%   0 out of 500

    ✓ status is 200 or 409

    CUSTOM
    booked_success.................: 16     142.217718/s

    HTTP
    http_req_duration..............: avg=41.5ms  min=8.92ms  med=38.09ms max=101.22ms p(90)=59.03ms p(95)=66.48ms
      { expected_response:true }...: avg=77.54ms min=59.62ms med=75.68ms max=101.22ms p(90)=99.56ms p(95)=100.32ms
    http_req_failed................: 96.80% 484 out of 500
    http_reqs......................: 500    4444.303688/s

    EXECUTION
    iteration_duration.............: avg=50.53ms min=10.43ms med=49.69ms max=107.88ms p(90)=69.57ms p(95)=78.98ms
    iterations.....................: 500    4444.303688/s

    NETWORK
    data_received..................: 85 kB  755 kB/s
    data_sent......................: 73 kB  647 kB/s




running (0m00.1s), 000/500 VUs, 500 complete and 0 interrupted iterations
hotspot ✓ [======================================] 500 VUs  00.1s/30s  500/500 shared iters

3) Verification
{
  "showId": "6a14654ba90b3b85d18e30d7",
  "seatCount": 1,
  "totalBookings": 16,
  "unbookedSeats": 0,
  "bookedSeats": 0,
  "oversoldCount": 1,
  "isValid": false,
  "oversoldSeats": [
    {
      "seatId": "6a14654ba90b3b85d18e30d8",
      "number": 1,
      "bookingCount": 16
    }
  ]
}
```
