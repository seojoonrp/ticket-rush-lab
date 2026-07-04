# Benchmark Summary

각 단계에서 동일 조건으로 **5회 반복 측정했을 때 median**이다.
재현 스크립트와 run별 원자료는 [`bench/`](../../bench) (`bench/out/raw.csv`)에서 확인할 수 있다.

## 측정 조건

- 환경: WSL2, Go 1.x, MongoDB 7 (Single node replica set, Docker) + Redis 7
- 부하 도구: k6
- Hotspot: 좌석 1개, 동시 요청 500 (shared-iterations, VU 500)
- Spread: 좌석 100개, 최대 500 VU, ramping (10s↑ / 20s 유지 / 5s↓)
- 집계: 단계당 5회 median, 서버 재기동 직후 1회는 워밍업으로 제외

## Hotspot

| 단계           | oversold 좌석 수 | 최대 중복 | isValid |
| -------------- | ---------------- | --------- | ------- |
| 1. Naive       | 1                | 18        | false   |
| 2. DB Atomic   | 0                | -         | true    |
| 3. Redis       | 0                | -         | true    |
| 4. Async Write | 0                | -         | true    |

> Hotspot은 latency를 제외했다. 총 500요청이 약 0.1초 만에 끝나 표본이 작아 p95가 노이즈에 지배되기 때문이다. 실제로 단계 간 median 차이보다 run 간 편차가 더 크다. Hotspot에서 볼 것은 정합성이다.

## Spread

| 단계           | 총 요청 | 처리량 (req/s) | p95 (total) | p95 (200 OK) | oversold 좌석 수 |
| -------------- | ------- | -------------- | ----------- | ------------ | ---------------- |
| 1. Naive       | 341,682 | 10,262         | 83.2ms      | 8.8ms        | 1                |
| 2. DB Atomic   | 336,825 | 10,116         | 84.1ms      | 9.2ms        | 0                |
| 3. Redis       | 508,549 | 15,480         | 51.5ms      | 9.7ms        | 0                |
| 4. Async Write | 517,581 | 15,559         | 51.7ms      | 1.7ms        | 0                |

## 요약

각 단계별로 바꾼 것과 trade-off 정리.

- **1. Naive** - race condition 발생. hotspot에서 한 좌석이 최대 중복 median 18건(관측 최대 67건)으로 oversell, spread에서도 소량(좌석 1개) oversell.
- **2. DB Atomic** - filter 조건을 건 atomic update로 oversell 차단 (hotspot, spread 모두 0). 처리량/latency는 spread에서 naive와 사실상 동일. 게이트가 여전히 Mongo write path라 병목이 그대로다.
- **3. Redis** - 좌석 점유 판정을 Redis로 이전. spread 처리량이 약 10k에서 15.5k req/s로, p95가 83에서 51ms로 뚜렷하게 향상. DB write는 always update로 되돌렸으나 게이트가 Redis라 정합성 유지.
- **4. Asynchronous Write** - 승자 판정 후 write를 worker pool로 넘겨 비동기 처리. 게이트가 Redis라 oversold 0 유지. 승자(200) latency가 spread에서 9.7에서 1.7ms로 크게 개선됐으나, 부하 대부분이 패자(409)라 total 처리량/p95(51.7ms)는 3단계와 사실상 동일.
