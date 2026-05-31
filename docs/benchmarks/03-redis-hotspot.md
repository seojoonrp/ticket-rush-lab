```bash
1) Register show
  show_id = 6a1c41c2bc1c327a59f18891
  seat_id = 6a1c41c2bc1c327a59f18892

2) Execute 500 concurrent requests

         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/


     execution: local
        script: /home/seojoonrp/workspace/study/ticket-rush-lab/loadtest/hotspot.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m0s max duration (incl. graceful stop):
              * hotspot: 500 iterations shared among 500 VUs (maxDuration: 30s, gracefulStop: 30s)



  █ TOTAL RESULTS

    checks_total.......: 500     2186.57818/s
    checks_succeeded...: 100.00% 500 out of 500
    checks_failed......: 0.00%   0 out of 500

    ✓ status is 200 or 409

    CUSTOM
    booked_success.................: 1      4.373156/s

    HTTP
    http_req_duration..............: avg=90.52ms min=3.93ms  med=84.75ms max=180.44ms p(90)=131.72ms p(95)=160.08ms
      { expected_response:true }...: avg=68.13ms min=68.13ms med=68.13ms max=68.13ms  p(90)=68.13ms  p(95)=68.13ms
    http_req_failed................: 99.80% 499 out of 500
    http_reqs......................: 500    2186.57818/s

    EXECUTION
    iteration_duration.............: avg=94.02ms min=4.01ms  med=87.25ms max=181.7ms  p(90)=133.88ms p(95)=162.9ms
    iterations.....................: 500    2186.57818/s

    NETWORK
    data_received..................: 86 kB  378 kB/s
    data_sent......................: 73 kB  319 kB/s




running (0m00.2s), 000/500 VUs, 500 complete and 0 interrupted iterations
hotspot ✓ [======================================] 500 VUs  00.2s/30s  500/500 shared iters

3) Verification
{
  "showId": "6a1c41c2bc1c327a59f18891",
  "seatCount": 1,
  "totalBookings": 1,
  "unbookedSeats": 0,
  "bookedSeats": 1,
  "oversoldCount": 0,
  "isValid": true,
  "oversoldSeats": null
}
```
