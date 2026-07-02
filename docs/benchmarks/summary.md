# Benchmark Summary

모든 측정은 동일 조건에서 수행한다. 단계별 raw 데이터와 상세는 각 단계 파일(`01-naive.md` 등)에서 확인할 수 있다.

## 측정 조건

- 환경: WSL2, Go 1.x, MongoDB 7 (Single node replica set, Docker)
- 부하 도구: k6
- Hotspot: 좌석 1개, 동시 요청 500 (shared-iterations, VU 500)
- Spread: 좌석 100개, 최대 500 VU, ramping (10s↑ / 20s 유지 / 5s↓)

## Hotspot

| 단계           | oversold 좌석 수 | 최대 중복 | isValid | p95 latency (total) |
| -------------- | ---------------- | --------- | ------- | ------------------- |
| 1. Naive       | 1                | 16        | false   | 66ms                |
| 2. DB Atomic   | 0                | -         | true    | 113ms               |
| 3. Redis       | 0                | -         | true    | 160ms               |
| 4. Worker pool | 0                | -         | true    | 155ms               |

## Spread

| 단계           | 총 요청 | 처리량(req/s) | p95 (total) | p95 (성공 200) | oversold 좌석 수 |
| -------------- | ------- | ------------- | ----------- | -------------- | ---------------- |
| 1. Naive       | 267,368 | 8,100         | 111ms       | 10.7ms         | 2                |
| 2. DB Atomic   | 265,808 | 7,944         | 112ms       | 11.45ms        | 0                |
| 3. Redis       | 431,557 | 13,009        | 64ms        | 12.68ms        | 0                |
| 4. Worker pool | 391,873 | 11,883        | 70.99ms     | 3.22ms         | 0                |

## 요약

각 단계별로 바꾼 것과 trade-off 정리.

- **1. Naive** — race condition 그대로. oversell 발생
- **2. DB Atomic** — filter 조건 추가. oversell 차단됨. spread에서는 latency가 비슷했으나 hotspot에서 p95 latency가 늘어나는 trade-off 발생
- **3. Redis** — 좌석 점유 판정을 Redis로 처리. spread에서의 처리량 및 p95 latency 향상. DB write는 무조건 update로 되돌렸으나 게이트가 Redis라 oversell 0 유지.
- **4. Worker pool** — 승자 판정 후 write를 worker pool로 넘겨 비동기 처리. 게이트가 Redis라 oversold 0 유지되고 승자 200 latency가 3.22ms로 크게 개선됐으나, 부하의 대부분이 패자(409)라 total 처리량/p95는 3단계와 노이즈 범위 안에서 겹친다.
