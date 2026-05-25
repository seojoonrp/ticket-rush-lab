```bash
1) Register show with 100 seats
./loadtest/reproduce_spread.sh: line 17: SEAT_ID: unbound variable
 seojoonrp@DESKTOP-CSJ ~/..../ticket-rush-lab  main  ./loadtest/reproduce_spread.sh     1) Register show with 100 seats
  show_id = 6a146588a90b3b85d18e314e

2) Execute requests on 100 different seats (max 500 VU)

         /\      Grafana   /‾‾/
    /\  /  \     |\  __   /  /
   /  \/    \    | |/ /  /   ‾‾\
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/


     execution: local
        script: /home/seojoonrp/practice/api_servers/ticket-rush-lab/loadtest/spread.js
        output: -

     scenarios: (100.00%) 1 scenario, 500 max VUs, 1m5s max duration (incl. graceful stop):
              * spread: Up to 500 looping VUs for 35s over 3 stages (gracefulRampDown: 30s, gracefulStop: 30s)



  █ THRESHOLDS

    http_req_duration
    ✓ 'p(95)<500' p(95)=111.52ms


  █ TOTAL RESULTS

    checks_total.......: 267368  8110.997527/s
    checks_succeeded...: 100.00% 267368 out of 267368
    checks_failed......: 0.00%   0 out of 267368

    ✓ status is 200 or 409

    CUSTOM
    booked_success.................: 102    3.094318/s

    HTTP
    http_req_duration..............: avg=49.66ms min=-1211567340ns med=45.21ms max=293.22ms p(90)=92.22ms p(95)=111.52ms
      { expected_response:true }...: avg=7.89ms  min=4.87ms        med=7.87ms  max=12.2ms   p(90)=9.99ms  p(95)=10.7ms
    http_req_failed................: 99.96% 267266 out of 267368
    http_reqs......................: 267368 8110.997527/s

    EXECUTION
    iteration_duration.............: avg=51.32ms min=1.08ms        med=46.65ms max=337.46ms p(90)=94.65ms p(95)=114.48ms
    iterations.....................: 267368 8110.997527/s
    vus............................: 6      min=6                max=500
    vus_max........................: 500    min=500              max=500

    NETWORK
    data_received..................: 46 MB  1.4 MB/s
    data_sent......................: 40 MB  1.2 MB/s




running (0m33.0s), 000/500 VUs, 267368 complete and 0 interrupted iterations
spread ✓ [======================================] 000/500 VUs  35s

3) Verification
{
  "showId": "6a146588a90b3b85d18e314e",
  "seatCount": 100,
  "totalBookings": 102,
  "unbookedSeats": 0,
  "bookedSeats": 98,
  "oversoldCount": 2,
  "isValid": false,
  "oversoldSeats": [
    {
      "seatId": "6a146588a90b3b85d18e3166",
      "number": 24,
      "bookingCount": 2
    },
    {
      "seatId": "6a146588a90b3b85d18e31aa",
      "number": 92,
      "bookingCount": 2
    }
  ]
}
```
