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
| 3. Redis       |                  |           |         |                     |
| 4. Worker pool |                  |           |         |                     |

## Spread

| 단계           | 총 요청 | 처리량(req/s) | p95 latency | oversold 좌석 수 |
| -------------- | ------- | ------------- | ----------- | ---------------- |
| 1. Naive       | 267,368 | 8,100         | 111ms       | 2                |
| 2. DB Atomic   | 265,808 | 7,944         | 112ms       | 0                |
| 3. Redis       |         |               |             |                  |
| 4. Worker pool |         |               |             |                  |

## 요약

각 단계별로 바꾼 것과 trade-off 정리.

- **1. Naive** — race condition 그대로. oversell 발생
- **2. DB Atomic** —
- **3. Redis** —
- **4. Worker pool** —
