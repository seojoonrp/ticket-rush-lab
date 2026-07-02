```bash
1) Register show with 100 seats
  show_id = 6a4647258c3fae4898c3945a

2) Execute requests on 100 different seats (max 500 VU)

         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/


     execution: local
        script: /home/seojoonrp/workspace/study/ticket-rush-lab/loadtest/spread.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m5s max duration (incl. graceful stop):
              * spread: Up to 500 looping VUs for 35s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)



  █ THRESHOLDS

    http_req_duration
    ✓ 'p(95)<500' p(95)=70.99ms


  █ TOTAL RESULTS

    checks_total.......: 391873  11883.871561/s
    checks_succeeded...: 100.00% 391873 out of 391873
    checks_failed......: 0.00%   0 out of 391873

    ✓ status is 200 or 409

    CUSTOM
    booked_success.................: 100    3.032582/s

    HTTP
    http_req_duration..............: avg=33.74ms min=-2361722905ns med=29.65ms max=702.86ms p(90)=57.75ms p(95)=70.99ms
      { expected_response:true }...: avg=2.03ms  min=1.03ms        med=1.9ms   max=4.18ms   p(90)=2.94ms  p(95)=3.22ms
    http_req_failed................: 99.97% 391773 out of 391873
    http_reqs......................: 391873 11883.871561/s

    EXECUTION
    iteration_duration.............: avg=34.92ms min=969.97µs      med=30.88ms max=332.56ms p(90)=60.65ms p(95)=76ms
    iterations.....................: 391873 11883.871561/s
    vus............................: 7      min=7                max=500
    vus_max........................: 500    min=500              max=500

    NETWORK
    data_received..................: 68 MB  2.1 MB/s
    data_sent......................: 59 MB  1.8 MB/s




running (0m33.0s), 000/500 VUs, 391873 complete and 0 interrupted iterations
spread ✓ [======================================] 000/500 VUs  35s

Draining write queue...

3) Verification
{
  "showId": "6a4647258c3fae4898c3945a",
  "seatCount": 100,
  "totalBookings": 100,
  "unbookedSeats": 0,
  "bookedSeats": 100,
  "oversoldCount": 0,
  "isValid": true,
  "oversoldSeats": null
}
```
