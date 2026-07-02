```bash
1) Register show with 100 seats
show_id = 6a15bf5caf5d8c56f063553c

2.  Execute requests on 100 different seats (max 500 VU)

            /\      Grafana   /‾‾/

    /\ / \ |\ ** / /
    / \/ \ | |/ / / ‾‾\
     / \ | ( | (‾) |
    / **\_\_\_\_**** \ |\_|\_\ \_\_\_\_\_/

        execution: local
           script: /home/seojoonrp/workspace/study/ticket-rush-lab/loadtest/spread.js
           output: -

        scenarios: (100.00%) 1 scenario, 500 max VUs, 1m5s max duration (incl. graceful stop):
                 * spread: Up to 500 looping VUs for 35s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)

█ THRESHOLDS

    http_req_duration
    ✓ 'p(95)<500' p(95)=112.4ms

█ TOTAL RESULTS

    checks_total.......: 265808  7944.205468/s
    checks_succeeded...: 100.00% 265808 out of 265808
    checks_failed......: 0.00%   0 out of 265808

    ✓ status is 200 or 409

    CUSTOM
    booked_success.................: 100    2.988701/s

    HTTP
    http_req_duration..............: avg=49.79ms min=-906402028ns med=44.96ms max=324.01ms p(90)=91.76ms p(95)=112.4ms
      { expected_response:true }...: avg=8ms     min=5.22ms       med=7.87ms  max=14.07ms  p(90)=10.17ms p(95)=11.45ms
    http_req_failed................: 99.96% 265708 out of 265808
    http_reqs......................: 265808 7944.205468/s

    EXECUTION
    iteration_duration.............: avg=51.58ms min=1.12ms       med=46.57ms max=324.1ms  p(90)=94.56ms p(95)=116.87ms
    iterations.....................: 265808 7944.205468/s
    vus............................: 5      min=5                max=500
    vus_max........................: 500    min=500              max=500

    NETWORK
    data_received..................: 46 MB  1.4 MB/s
    data_sent......................: 40 MB  1.2 MB/s

running (0m33.5s), 000/500 VUs, 265808 complete and 0 interrupted iterations
spread ✓ [======================================] 000/500 VUs 35s

3. Verification
   {
   "showId": "6a15bf5caf5d8c56f063553c",
   "seatCount": 100,
   "totalBookings": 100,
   "unbookedSeats": 0,
   "bookedSeats": 100,
   "oversoldCount": 0,
   "isValid": true,
   "oversoldSeats": null
   }

```

```

```
