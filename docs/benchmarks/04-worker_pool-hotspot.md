```bash
1) Register show
  show_id = 6a46471f8c3fae4898c39457
  seat_id = 6a46471f8c3fae4898c39458

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

    checks_total.......: 500     2418.000548/s
    checks_succeeded...: 100.00% 500 out of 500
    checks_failed......: 0.00%   0 out of 500

    ✓ status is 200 or 409

    CUSTOM
    booked_success.................: 1      4.836001/s

    HTTP
    http_req_duration..............: avg=86.9ms  min=30.15ms med=80.58ms max=191.62ms p(90)=150.06ms p(95)=155.68ms
      { expected_response:true }...: avg=32.57ms min=32.57ms med=32.57ms max=32.57ms  p(90)=32.57ms  p(95)=32.57ms
    http_req_failed................: 99.80% 499 out of 500
    http_reqs......................: 500    2418.000548/s

    EXECUTION
    iteration_duration.............: avg=98.65ms min=36.76ms med=96ms    max=203.39ms p(90)=159.44ms p(95)=167.1ms
    iterations.....................: 500    2418.000548/s

    NETWORK
    data_received..................: 86 kB  418 kB/s
    data_sent......................: 73 kB  353 kB/s




running (0m00.2s), 000/500 VUs, 500 complete and 0 interrupted iterations
hotspot ✓ [======================================] 500 VUs  00.2s/30s  500/500 shared iters

Draining write queue...

3) Verification
{
  "showId": "6a46471f8c3fae4898c39457",
  "seatCount": 1,
  "totalBookings": 1,
  "unbookedSeats": 0,
  "bookedSeats": 1,
  "oversoldCount": 0,
  "isValid": true,
  "oversoldSeats": null
}
```
