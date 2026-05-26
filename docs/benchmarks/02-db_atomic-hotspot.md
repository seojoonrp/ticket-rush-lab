```bash
1) Register show
  show_id = 6a15bf53af5d8c56f0635539
  seat_id = 6a15bf54af5d8c56f063553a

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

    checks_total.......: 500     3103.647275/s
    checks_succeeded...: 100.00% 500 out of 500
    checks_failed......: 0.00%   0 out of 500

    ✓ status is 200 or 409

    CUSTOM
    booked_success.................: 1      6.207295/s

    HTTP
    http_req_duration..............: avg=88.33ms  min=26.42ms  med=97.17ms  max=128.35ms p(90)=109.37ms p(95)=113.54ms
      { expected_response:true }...: avg=103.41ms min=103.41ms med=103.41ms max=103.41ms p(90)=103.41ms p(95)=103.41ms
    http_req_failed................: 99.80% 499 out of 500
    http_reqs......................: 500    3103.647275/s

    EXECUTION
    iteration_duration.............: avg=94.34ms  min=35.04ms  med=100.83ms max=149.29ms p(90)=114.74ms p(95)=118.51ms
    iterations.....................: 500    3103.647275/s

    NETWORK
    data_received..................: 86 kB  536 kB/s
    data_sent......................: 73 kB  452 kB/s




running (0m00.2s), 000/500 VUs, 500 complete and 0 interrupted iterations
hotspot ✓ [======================================] 500 VUs  00.2s/30s  500/500 shared iters

3) Verification
{
  "showId": "6a15bf53af5d8c56f0635539",
  "seatCount": 1,
  "totalBookings": 1,
  "unbookedSeats": 0,
  "bookedSeats": 1,
  "oversoldCount": 0,
  "isValid": true,
  "oversoldSeats": null
}
```
