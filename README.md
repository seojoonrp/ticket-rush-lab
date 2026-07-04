# ticket-rush-lab

그동안 만든 프로젝트들은 트래픽이 몰릴 일이 없었다. 적당히 구현해도 성능이나 정합성 문제가 드러나지 않는 환경이었기에 최적화를 해도 티가 안 나는 경우가 대부분이었다. 게다가 주력 언어로 Go를 쓰면서도 정작 그 강점인 동시성을 제대로 써본 적이 없다는 것도 항상 마음에 걸렸다. 컴퓨터조직론 과목에서 마침 cache consistency를 재밌게 배우며 동시성에 흥미가 생긴 김에, 고동시성·고트래픽 환경에 직접 부딪혀보고 싶어 시작한 프로젝트다. Go를 주력으로 쓰는 개발자가 동시성이 뭔지 모르는게 말이 되는가?

주제는 티켓팅 서버. 한정된 좌석을 여러 사람이 동시에 예매하려고 몰려들 때 생기는 문제들을 직접 재현하고, 단계별로 해결해 나가면서 그 전후 변화를 기록한다.

처음부터 완벽한 티켓팅 서버를 만드는 게 아니라, **동시성 문제를 고려하지 않고 naive하게 구현한 후 생기는 문제점을 파악하고, 이를 단계적으로 정합성/성능 면에서 개선해나가며 그 과정에서의 벤치마크 변화를 측정해 기록하는 것**이다.

## 다루는 문제

티켓팅의 본질적인 어려움은 재고가 한정되어 있는데 수많은 요청이 동시에 들어온다는 데 있다. 좌석 하나를 두 사람이 거의 동시에 잡으려 하면 어떻게 처리해야 하는가? 한 좌석에 둘 이상의 사람들이 동시에 예매 요청을 날리면 naive하게 구현된 서버에서는 둘 다 예매되었다고 처리되는 oversell 문제가 발생한다. 근데 그렇다고 한명 한명 줄 세워서 처리하기에는 인기 콘서트는 수만 명이 동시에 요청을 날릴텐데... 너무 느리다. 문제 해결이 필요하다.

이 문제는 두 가지 관점에서 개선될 수 있고, 본 프로젝트는 이를 순차적으로 수행한다.

- **정합성**: oversell이 일어나는가. 한 좌석이 정확히 한 명에게만 팔리는가.
- **성능**: 정합성을 유지하면서 얼마나 많은 요청을, 얼마나 빠르게 처리하는가.

## 진행 단계

1. **Naive implementation** - Race condition을 해결하지 않은 naive한 구현으로 oversell을 재현한다.
2. **DB-level atomicity** - MongoDB의 단일 문서 atomic 연산으로 oversell을 막는다.
3. **Redis** - 재고 차감 처리를 메모리로 끌어올려 성능을 향상한다.
4. **Asynchronous write** - 무거운 write 작업을 worker pool로 비동기 처리한다.

각 단계마다 같은 부하를 주고 정확성과 성능이 어떻게 바뀌는지 기록한다.

## 벤치마크 Summary

각 단계에서 동일 조건으로 **5회 반복 측정했을 때 median**이다.
Raw data는`bench/out/raw.csv`에서 확인할 수 있다.

**측정 조건**

- 환경: WSL2, Go 1.x, MongoDB 7 (Single node replica set, Docker) + Redis 7
- 부하 도구: k6
- Hotspot: 좌석 1개에 동시 요청 500 (shared-iterations, VU 500). 모두가 같은 좌석을 노리므로 정합성 검증에 효과적이다.
- Spread: 좌석 100개에 최대 500 VU ramping (10s↑ / 20s 유지 / 5s↓). 경합이 흩어져 현실 티켓팅에 더 가깝고 성능 분석에 쓴다.

### Hotspot

| 단계           | oversold 좌석 수 | 최대 중복 | isValid |
| -------------- | ---------------- | --------- | ------- |
| 1. Naive       | 1                | 18        | false   |
| 2. DB Atomic   | 0                | -         | true    |
| 3. Redis       | 0                | -         | true    |
| 4. Async write | 0                | -         | true    |

> Hotspot은 latency를 제외했다. 총 500요청이 약 0.1초 만에 끝나 표본이 작아 p95가 노이즈에 지배되기 때문이다. 실제로 단계 간 median 차이보다 run 간 편차가 더 크다. Hotspot에서 볼 것은 정합성이다.

### Spread

| 단계           | 총 요청 | 처리량 (req/s) | p95 (total) | p95 (200 OK) | oversold 좌석 수 |
| -------------- | ------- | -------------- | ----------- | ------------ | ---------------- |
| 1. Naive       | 341,682 | 10,262         | 83.2ms      | 8.8ms        | 1                |
| 2. DB Atomic   | 336,825 | 10,116         | 84.1ms      | 9.2ms        | 0                |
| 3. Redis       | 508,549 | 15,480         | 51.5ms      | 9.7ms        | 0                |
| 4. Async write | 517,581 | 15,559         | 51.7ms      | 1.7ms        | 0                |

단계별 상세 분석은 아래 [Naive](#1-naive-implementation-260526) ~ [Async write](#4-worker-pool-260702) 섹션에서 다룬다.

## 프로젝트 구조

```
cmd/server          진입점, DI
internal/
  handler           HTTP 핸들러
  service           예매/검증 로직. 단계별로 바뀌는 부분이 여기 모여 있다.
  repository        MongoDB 처리
  middleware        X-User-ID 검증 미들웨어
  model             도메인 타입
  apperr            에러 타입과 중앙 error handler
loadtest            k6 부하 스크립트와 재현용 shell 스크립트
bench               벤치마크 재현 harness
```

Handler - Service - Repository - MongoDB의 layered 구조에 Redis와 worker pool을 얹었다. 의존성은 전부 인터페이스로 두고 프레임워크 없이 `main.go`에서 직접 조립한다. 좌석 점유를 판정하고 예매를 확정하는 로직은 전부 Service의 `Book()`에 모여 있어, 네 단계에 걸친 개선이 이 `Book()` 안에서만 일어나고 나머지 계층은 거의 그대로 유지된다. Repository는 저장소별로 나눠 `seat`·`booking`·`show`는 MongoDB를, `seat_claim`은 Redis를 담당한다.

## API

| Method | Path                | Description                                                                    |
| ------ | ------------------- | ------------------------------------------------------------------------------ |
| POST   | `/shows`            | 공연과 좌석을 생성한다. 매 측정마다 새 공연을 만들어 깨끗한 상태에서 시작한다. |
| POST   | `/seats/:id/book`   | 좌석을 예매한다. 사용자 식별은 `X-User-ID` 헤더로 받는다.                      |
| GET    | `/shows/:id/verify` | 좌석별 booking을 집계해 oversell을 적발한다.                                   |

예매 성공은 200, 이미 찬 좌석은 409로 응답한다. 매진으로 인한 거절(409)과 서버 오류(500)를 구분해야 부하 측정에서 정상적인 경합과 의도되지 않은 오류들이 섞이지 않는다.

## 실행

MongoDB와 Redis는 Docker로, 서버는 로컬에서 띄운다. 접속 정보는 `.env`로 주입한다(`.env.example` 참고).

```
cp .env.example .env
docker compose up -d
go run ./cmd/server/main.go
```

단일 부하 재현은 `loadtest`의 shell 스크립트로 한다. 실행 중인 서버에 한 번 부하를 주고 결과를 출력한다.

```
./loadtest/reproduce_hotspot.sh    # hotspot (한 좌석에 집중 요청 - 500개의 요청이 동시에 들어오는 시나리오)
./loadtest/reproduce_spread.sh     # spread (여러 좌석에 분산 - 100개 좌석에 35초간 최대 500 VU로 계속해서 요청)
```

위 [벤치마크 결과 요약](#벤치마크-결과-요약)의 median 값은 `bench/run.sh`로 재현한다. 각 단계 커밋을 자동으로 오가며 서버를 빌드·실행하고, 단계별로 5회씩 측정해 `bench/out`에 집계한다.

```
./bench/run.sh          # 전체 단계
./bench/run.sh 3-redis  # 특정 단계만
```

> 측정 환경은 WSL2다. k6 지연 측정에서 시계 점프로 인한 음수 시간이 일부 관찰됐는데, 중앙값/p95/처리량 같은 핵심 지표는 영향을 받지 않아 그대로 사용한다.

<br>

## 1. Naive implementation (260526)

Naive한 예매 로직은 아래와 같은 단계로 수행된다.

1. 좌석을 읽고
2. 비었는지 검사하고
3. 점유 표시

문제는 검사와 점유 사이에 틈이 있다는 것이다. 한 요청이 비어있다고 판단한 직후, 점유를 기록하기 전에 다른 요청이 끼어들면 그 요청도 비어있다고 판단하는 race condition이 생긴다. 결국 둘 다 통과해서 같은 좌석에 booking이 두 건 박힌다.

**Hotspot**

```
seatCount: 1
oversoldCount: 1
oversoldSeats: [{ number: 1, bookingCount: 16 }]
```

좌석 하나에 16건의 예매가 기록되는 기념비적인 oversell이 일어났다. 484건은 정상적으로 409 Conflict를 받았지만, 16건이 동시에 검사를 통과해 버린 것이다.

**Spread**

```
seatCount: 100
http_reqs: 267,368 (약 8,100 req/s)
http_req_duration p95: 111ms
oversoldCount: 2 (좌석 2개가 각각 2건 중복)
```

경합을 100개 좌석으로 분산시키자 같은 race condition인데도 oversell이 훨씬 드물게 일어났다. Hotspot에서는 요청을 500개밖에 안 날렸는데도 16개의 요청이 중복됐지만, spread에서는 26만 건이 넘는 요청 중 단 2개 좌석만 2번씩 중복됐다. 경합이 한 점에 몰리느냐 흩어지느냐에 따라 발현 빈도가 크게 달라진다는 걸 보여준다.

이 수치는 이후 단계의 기준선이다. 2단계에서 원자적 연산을 적용하면 oversell이 0이 되어야 하고, 그때 처리량과 지연이 어떻게 변하는지가 다음 비교 대상이다.

<br>

## 2. DB-level atomic implementation (260527)

1단계의 문제는 검사와 점유가 두 번의 DB 왕복으로 쪼개져 있다는 것이었다. 그 사이의 틈이 race condition의 발생 지점이다. 그렇다면 틈을 없애려면 검사와 점유를 하나의 연산으로 합치면 된다.

방법은 의외로 단순하다. 검사 조건을 update의 필터 안으로 밀어넣는다.

```go
filter := bson.M{
  "_id":    id,
  "status": model.SeatAvailable,
}
```

**Hotspot**

```
seatCount: 1
booked_success: 1
oversoldCount: 0
isValid: true
```

500개의 동시 요청 중 정확히 1건만 좌석을 잡고 나머지 499건은 의도된대로 409를 받았다. oversell 방지 성공!

**Spread**

```
seatCount: 100
http_reqs: 265,808 (약 7,944 req/s)
http_req_duration p95: 112ms
oversoldCount: 0
```

100개 좌석에 분산시킨 26만 건의 요청에서도 oversold가 사라졌다. 좌석 100개에 정확히 100건의 booking이 박혔다.

흥미로운 건 성능이다. 정확성을 얻으면 성능을 잃는다는 trade-off의 직관과 달리, spread의 처리량과 지연은 1단계와 거의 동일했다(8,100 - 7,944 req/s, p95 111 - 112ms). 사실상 측정 노이즈 범위 안이라고 할 수 있다. 분산 부하에서는 좌석마다 경합이 흩어져 있어 조건부 update의 직렬화 비용이 잘 드러나지 않기 때문이다. 정확성을 사실상 공짜로 얻었다.

그런데 hotspot은 다른 그림을 보여준다. 한 좌석에 모든 요청이 몰리는 상황에서는 p95 latency가 66ms에서 113ms로 늘었다. 1단계의 update는 조건이 없어서 서로를 막지 않고 그냥 덮어쓰고 지나갔지만, 2단계의 조건부 update는 같은 문서를 두고 진짜로 경합을 해버린다. 그 비용이 한 점에 몰릴 때 지연으로 드러난 것이라 추정된다. 물론 hotspot은 요청 수가 워낙 적기 때문에 이를 정량적으로 분석하는 데는 좀 무리가 있다. 그래도 확실하게 지연이 생긴 건 알아볼 수 있다.

결론적으로 정확성을 위해 도입한 atomicity로 인해 hotspot 시나리오에서 약간의 지연이 생기는 trade-off가 있었다. 이 비용을 어떻게 줄일 것인가가 다음 단계의 출발점이다. 3단계에서는 redis를 도입해 재고 차감을 메모리로 끌어올려, 매진된 좌석을 향한 요청이 DB까지 내려가기 전에 빠르게 처리되도록 한다.

> **3단계 구현 이후의 회고... (260531)** <br> Hotspot에서의 p95 latency는 크게 의미가 없다. 요청을 500개밖에 안 날려서 표본이 너무 적기 때문이다. 작은 노이즈에도 결과가 너무 달라지기 때문에 좋은 성능 지표가 될 수 없고, 따라서 이를 trade-off라 말하기도 애매하다. 실제로 Redis를 적용한 3단계 결과를 보면 spread의 성능은 확실하게 올라가지만 hotspot의 p95 latency는 오히려 늘어난다.

<br>

## 3. Redis (260531)

2단계는 정확성을 잡았지만 비용을 남겼다. hotspot에서 조건부 update의 직렬화로 지연이 늘었고, 더 근본적으로는 spread에서 매진 이후 들어오는 요청 - 26만 건의 99.97% - 이 전부 DB까지 내려가 조건부 update를 날리고 이선좌(MatchedCount 0) 에러를 받아갔다. 티켓팅 서버의 진짜 부하는 표를 잡는 소수가 아니라, 매진인데도 새로고침하며 들이치는 다수다. 그 다수를 DB까지 보내지 않고 Redis를 사용해 메모리에서 거르는 게 이 단계의 목표다.

Redis는 명령 실행이 단일 스레드라 개별 명령이 atomic하게 처리된다. 즉 `SET key NX`는 키가 없을 때만 set하고 성공 여부를 돌려주는 atomic test-and-set이라, 기존 로직을 그대로 대체할 수 있다.

```go
// Seat Claim Repository
func (r *SeatClaimRepo) Claim(ctx context.Context, seatID primitive.ObjectID, userID string) (bool, error) {
	key := "seat:" + seatID.Hex()
	// value를 userID로 설정해 디버깅이 편하도록 (누가 이 좌석을 선점했는지 파악 가능)
	// expiration을 0으로 설정하면 만료되지 않음
	return r.rdb.SetNX(ctx, key, userID, 0).Result()
}

// seat:<id> 키를 NX로 선점. true면 예매 성공, false면 이선좌
claimed, _ := claimRepo.Claim(ctx, seatID, userID)
if !claimed {
    return apperr.ErrSeatTaken // 409, DB를 건드리지 않는다
}
```

예매 판정이 Redis로 넘어가면서 MongoDB는 동시성 방어선에서 내려와 그냥 저장소가 됐다. 그래서 1단계 이후로 쓸 일이 없을 것 같았던 UpdateOnBook 메서드를 부활시켜 다시 사용했다. 이래서 안 쓰는 코드 함부로 지우면 안되는 것 같다.

**Hotspot**

```
seatCount: 1
booked_success: 1
oversoldCount: 0
isValid: true
```

MongoDB에 저장하는 메서드는 1단계의 것을 그대로 썼는데, Redis에서 유효성 판단을 해주므로 정확성이 유지됐다.

**Spread**

```
seatCount: 100
http_reqs: 431,557 (약 13,009 req/s)
http_req_duration p95: 64ms
oversoldCount: 0
```

여기서 효과가 드러난다. 처리량이 7,944에서 13,009 req/s로 뛰고 p95가 112ms에서 64ms로 거의 반토막 났다. 이선좌 요청 43만 건이 빠른 Redis에서 잘려나가고 MongoDB는 좌석을 실제로 잡은 100명만 봤기 때문이다. 매진 트래픽을 메모리에서 거르는 것만으로 정확성을 유지한 채 성능까지 압도적으로 끌어올렸다.

hotspot의 p95 latency(160ms)는 2단계 회고에 적어둔 대로 오히려 전보다 늘었지만, 이는 표본 부족으로 인해 정량적 분석이 의미 없는 데이터다. Redis의 단일 키 방식이 가지는 본질적인 문제라 4단계에서도 딱히 좋아지지는 않을 거라 예상된다. hotspot은 최악의 경합 케이스로 남겨둔다.

> hotspot 지연이 go-redis 커넥션 풀 경합 때문인지 의심돼 풀을 기본값(GOMAXPROCS 확인 결과 120)에서 600으로 키워봤다. 하지만 표본이 안정적인 spread의 처리량이 움직이지 않아(13,009 to 12,281) 가설을 기각했다. 병목은 커넥션 풀이 아니라 요청 경로 자체에 있는 것 같다.

근데 아직 write가 동기다. 좌석 예매 성공 시 여전히 요청 안에서 seat update와 booking insert를 끝내고서야 응답을 받는다. 4단계에서는 이 무거운 write를 worker pool로 넘겨 비동기로 처리하고, 응답을 그만큼 앞당긴다.

<br>

## 4. Worker pool (260702)

3단계 끝에 남은 숙제는 write가 아직 동기라는 것이었다. 좌석을 잡은 요청은 Redis claim으로 이미 승부가 났는데도, seat update와 booking insert라는 Mongo write 두 방을 요청 안에서 끝내고서야 200을 받았다. 정확성 방어선은 Redis에 있으니 write는 요청 경로에서 떼어내 언제 하든 상관없어야 한다. 이걸 worker pool로 증명하는 게 4단계다.

승자 판정(SetNX)이 끝나면 job을 buffered channel에 넣고 곧바로 리턴한다. 채널 뒤에 붙은 워커 goroutine들이 Mongo write를 비동기로 처리한다.

```go
func (s *BookingService) Book(ctx context.Context, seatID primitive.ObjectID, userID string) error {
	claimed, _ := s.claimRepo.Claim(ctx, seatID, userID)
	if !claimed {
		return apperr.ErrSeatTaken
	}

	s.pool.Submit(Job{seatID: seatID, userID: userID})
	return nil // write를 안 기다리고 리턴
}
```

write가 비동기가 되면서 측정에 구멍이 하나 생긴다. verify는 booking을 집계해 정확성을 판정하는데, k6가 끝난 직후 큐에 write가 남아 있으면 booking이 실제보다 적게 잡혀 오탐이 난다. 그래서 큐 depth를 노출하는 `/health/queue`를 두고, reproduce 스크립트가 depth 0을 확인한 뒤에 verify를 때리도록 했다. 이건 새 성능 지표가 아니라 verify 타이밍만 보정하는 drain 확인용이라, 1\~3단계의 측정 기준(hotspot 정확성 + spread 처리량/p95)은 그대로 유지된다.

**Hotspot**

```
seatCount: 1
booked_success: 1
oversoldCount: 0
isValid: true
```

**Spread**

```
seatCount: 100
http_reqs: 391,873 (약 11,883 req/s)
http_req_duration p95: 70.99ms
totalBookings: 100
oversoldCount: 0
```

write를 요청 밖으로 뺐는데도 정확성이 그대로다. hotspot은 승자 1명에 oversold 0, spread는 좌석 100개가 정확히 다 팔리고 oversold 0, 둘 다 isValid: true. 방어선이 Redis에 있으니 write를 언제 하든 결과가 안 바뀐다. Drain polling 덕분에 spread에서도 정확성이 유지된다.

그런데 total 지표는 3단계와 사실상 구분이 안 된다. spread 처리량 11,883 req/s, total p95 70.99ms는 3단계(약 13,009 req/s, 64ms)와 노이즈 범위 안에서 겹치고 오히려 살짝 낮게도 나온다. 이유는 부하의 99.97%가 패자(409)라 승자 write를 아무리 비동기로 처리해도 total에는 영향을 미치지 않기 때문이다. 병목은 패자 경로(Redis 직렬화)에 있지 승자 write에 있지 않다. 4단계는 승자 경로만 건드렸는데 승자가 100명뿐이라 total 지표는 그 개선을 못 본다.

진짜 봐야 할 숫자는 처리량이 아니라 승자의 200 OK latency다. spread에서 `{ expected_response:true }`(200을 받은 승자들)의 p95는 1, 2, 3단계 각각 10.7, 11.45, 12.68ms로 10ms대에 묶여 있었다. 승자가 매 단계 Mongo write를 동기로 기다렸기 때문이다. 4단계에서 그 write를 worker pool로 빼자 3.22ms로 크게 감소했다. 승자가 write 세 방을 안 기다리고 claim 직후 리턴하니 200 응답이 이만큼 빨라진 것이다.
