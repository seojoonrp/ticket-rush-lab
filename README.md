# ticket-rush-lab

Go로 API 서버를 만들다 보니 goroutine이나 channel을 정작 제대로 써본 적이 없다는 걸 깨달았다. 컴퓨터조직론 과목에서 마침 cache consistency를 배우고 있어 동시성에 흥미가 생긴 김에 동시성을 제대로 공부해보고 싶어 시작한 프로젝트다. Go를 주력으로 쓰는 개발자가 동시성이 뭔지 모르는게 말이 되는가?

주제는 티켓팅 서버. 한정된 좌석을 여러 사람이 동시에 예매하려고 몰려들 때 생기는 문제들을 직접 재현하고, 단계별로 해결해 나가면서 그 전후 변화를 기록한다.

처음부터 완벽한 티켓팅 서버를 만드는 게 아니라, **동시성 문제를 일부러 터뜨린 다음 단계적으로 고치고, 개선해나가며 그 과정에서의 벤치마크 등의 변화를 측정해 기록하는 것**이다.

## 다루는 문제

티켓팅의 본질적인 어려움은 재고가 한정되어 있는데 수많은 요청이 동시에 들어온다는 데 있다. 좌석 하나를 두 사람이 거의 동시에 잡으려 하면 어떻게 처리해야 하는가? 한 좌석에 둘 이상의 사람들이 동시에 예매 요청을 날리면 naive하게 구현된 서버에서는 둘 다 예매되었다고 처리되는 oversell 문제가 발생한다. 근데 그렇다고 한명 한명 줄 세워서 처리하기에는 인기 콘서트는 수만 명이 동시에 요청을 날릴텐데... 너무 느리다. 문제 해결이 필요하다.

이 문제는 두 가지 관점에서 개선될 수 있고, 본 프로젝트는 이를 순차적으로 수행한다.

- **정확성**: oversell이 일어나는가. 한 좌석이 정확히 한 명에게만 팔리는가.
- **성능**: 같은 정확성을 유지하면서 얼마나 많은 요청을, 얼마나 빠르게 처리하는가.

## 진행 단계

1. **Naive implementation** - Race condition을 해결하지 않은 naive한 구현으로 oversell을 재현한다.
2. **DB-level atomicity** - MongoDB의 단일 문서 atomic 연산으로 oversell을 막는다.
3. **Redis** - 재고 차감 처리를 메모리로 끌어올려 성능을 향상한다.
4. **Worker pool** - 무거운 write 작업을 worker pool로 비동기 처리한다.

각 단계 사이에서 같은 부하를 주고 정확성과 성능이 어떻게 바뀌는지 기록한다.

## 구조

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
docs/benchmarks     단계별 측정 결과
```

좌석을 한 명만 차지할 수 있는 독립된 슬롯으로 보고, `seats`를 별도 컬렉션으로 뒀다. 이렇게 하면 좌석마다 경합이 독립적이라 "한 좌석에 몰리는 부하"와 "전체 트래픽 부하"를 나눠서 실험할 수 있다.

oversell을 감지하는 방식이 이 프로젝트의 한 축이다. 좌석의 상태 필드만 보면 나중에 쓴 사람이 앞사람을 덮어써서 중복의 흔적이 사라진다. 그래서 예매가 성공할 때마다 `bookings` 컬렉션에 레코드를 INSERT하고, 검증 단계에서 좌석별로 booking을 집계해 2건 이상인 좌석을 찾아낸다.

## API

| Method | Path                | Description                                                                    |
| ------ | ------------------- | ------------------------------------------------------------------------------ |
| POST   | `/shows`            | 공연과 좌석을 생성한다. 매 측정마다 새 공연을 만들어 깨끗한 상태에서 시작한다. |
| POST   | `/seats/:id/book`   | 좌석을 예매한다. 사용자 식별은 `X-User-ID` 헤더로 받는다.                      |
| GET    | `/shows/:id/verify` | 좌석별 booking을 집계해 oversell을 적발한다.                                   |

예매 성공은 200, 이미 찬 좌석은 409로 응답한다. 매진으로 인한 거절(409)과 서버 오류(500)를 구분해야 부하 측정에서 정상적인 경합과 의도되지 않은 오류들이 섞이지 않는다.

## 실행

MongoDB는 Docker로, 서버는 로컬에서 띄운다.

```
docker compose up -d
go run ./cmd/server/main.go
```

부하 재현은 `loadtest` 안의 shell 스크립트를 통해 수행한다.

```
./loadtest/reproduce_hotspot.sh    # hotspot (한 좌석에 집중 요청 - 500개의 요청이 동시에 들어오는 시나리오)
./loadtest/reproduce_spread.sh     # spread (여러 좌석에 분산 - 100개 좌석에 35초간 최대 500 VU로 계속해서 요청)
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
