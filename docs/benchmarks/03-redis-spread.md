```bash
1) Register show with 100 seats
  show_id = 6a1c41cdbc1c327a59f18894

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
    ✓ 'p(95)<500' p(95)=63.85ms


  █ TOTAL RESULTS

    checks_total.......: 431557  13009.535255/s
    checks_succeeded...: 100.00% 431557 out of 431557
    checks_failed......: 0.00%   0 out of 431557

    ✓ status is 200 or 409

    CUSTOM
    booked_success.................: 100    3.014558/s

    HTTP
    http_req_duration..............: avg=30.24ms min=-1048706075ns med=27.98ms max=204.24ms p(90)=53.41ms p(95)=63.85ms
      { expected_response:true }...: avg=9.34ms  min=5.52ms        med=9.31ms  max=14.19ms  p(90)=12.05ms p(95)=12.68ms
    http_req_failed................: 99.97% 431457 out of 431557
    http_reqs......................: 431557 13009.535255/s

    EXECUTION
    iteration_duration.............: avg=31.71ms min=938.4µs       med=29.13ms max=204.37ms p(90)=55.86ms p(95)=67.25ms
    iterations.....................: 431557 13009.535255/s
    vus............................: 5      min=5                max=500
    vus_max........................: 500    min=500              max=500

    NETWORK
    data_received..................: 75 MB  2.3 MB/s
    data_sent......................: 65 MB  1.9 MB/s




running (0m33.2s), 000/500 VUs, 431557 complete and 0 interrupted iterations
spread ✓ [======================================] 000/500 VUs  35s

3) Verification
{
  "showId": "6a1c41cdbc1c327a59f18894",
  "seatCount": 100,
  "totalBookings": 100,
  "unbookedSeats": 0,
  "bookedSeats": 100,
  "oversoldCount": 0,
  "isValid": true,
  "oversoldSeats": null
}
```
