#!/usr/bin/env python3
"""Turn per-run k6/verify JSON into raw.csv + median summary tables.
usage: aggregate.py <results_dir> <out_dir>   (run-0 = warmup, ignored)
"""
import csv, glob, json, os, re, statistics, sys

ORDER = ["1-naive", "2-db_atomic", "3-redis", "4-worker_pool"]
LABEL = {"1-naive": "1. Naive", "2-db_atomic": "2. DB Atomic",
         "3-redis": "3. Redis", "4-worker_pool": "4. Worker pool"}


def val(d, name, stat):
    try:
        return d["metrics"][name]["values"][stat]
    except (KeyError, TypeError):
        return None


def load(results_dir):
    rows = []
    for sj in sorted(glob.glob(os.path.join(results_dir, "*", "*", "run-*.summary.json"))):
        idx = int(re.search(r"run-(\d+)\.summary\.json$", sj).group(1))
        if idx == 0:                      # warmup
            continue
        scn = os.path.basename(os.path.dirname(sj))
        stage = os.path.basename(os.path.dirname(os.path.dirname(sj)))
        with open(sj) as f:
            d = json.load(f)
        v = {}
        vp = sj.replace(".summary.json", ".verify.json")
        if os.path.exists(vp):
            with open(vp) as f:
                v = json.load(f)
        dups = [s.get("bookingCount", 0) for s in (v.get("oversoldSeats") or [])]
        rows.append({
            "stage": stage, "scenario": scn, "run": idx,
            "p95_total_ms": val(d, "http_req_duration", "p(95)"),
            "p95_ok200_ms": val(d, "http_req_duration{expected_response:true}", "p(95)"),
            "http_reqs": val(d, "http_reqs", "count"),
            "req_per_s": val(d, "http_reqs", "rate"),
            "booked_success": val(d, "booked_success", "count"),
            "oversold_count": v.get("oversoldCount"),
            "max_dup": max(dups) if dups else 0,
            "total_bookings": v.get("totalBookings"),
        })
    return rows


def med(rows, key):
    xs = [r[key] for r in rows if r.get(key) is not None]
    return statistics.median(xs) if xs else None


def i(v): return "-" if v is None else str(int(round(v)))
def ms(v): return "-" if v is None else f"{round(v, 2)}ms"


def sub(rows, stage, scn):
    return [r for r in rows if r["stage"] == stage and r["scenario"] == scn]


def hotspot_md(rows):
    out = ["## Hotspot", "",
           "| 단계 | oversold 좌석 수 | 최대 중복 | isValid | p95 latency (total) |",
           "| ---- | ---- | ---- | ---- | ---- |"]
    for st in ORDER:
        g = sub(rows, st, "hotspot")
        if not g:
            continue
        ov = med(g, "oversold_count")
        dup = i(med(g, "max_dup")) if ov and ov > 0 else "-"
        out.append(f"| {LABEL[st]} | {i(ov)} | {dup} | "
                   f"{'true' if ov == 0 else 'false'} | {ms(med(g, 'p95_total_ms'))} |")
    return "\n".join(out) + "\n"


def spread_md(rows):
    out = ["## Spread", "",
           "| 단계 | 총 요청 | 처리량(req/s) | p95 (total) | p95 (성공 200) | oversold 좌석 수 |",
           "| ---- | ---- | ---- | ---- | ---- | ---- |"]
    for st in ORDER:
        g = sub(rows, st, "spread")
        if not g:
            continue
        out.append(f"| {LABEL[st]} | {i(med(g, 'http_reqs'))} | {i(med(g, 'req_per_s'))} | "
                   f"{ms(med(g, 'p95_total_ms'))} | {ms(med(g, 'p95_ok200_ms'))} | "
                   f"{i(med(g, 'oversold_count'))} |")
    return "\n".join(out) + "\n"


def main():
    results_dir, out_dir = sys.argv[1], sys.argv[2]
    os.makedirs(out_dir, exist_ok=True)
    rows = load(results_dir)
    if not rows:
        sys.exit("no measured runs found under " + results_dir)
    cols = ["stage", "scenario", "run", "p95_total_ms", "p95_ok200_ms", "http_reqs",
            "req_per_s", "booked_success", "oversold_count", "max_dup", "total_bookings"]
    with open(os.path.join(out_dir, "raw.csv"), "w", newline="") as f:
        w = csv.DictWriter(f, fieldnames=cols)
        w.writeheader()
        for r in sorted(rows, key=lambda r: (r["stage"], r["scenario"], r["run"])):
            w.writerow(r)
    h, s = hotspot_md(rows), spread_md(rows)
    open(os.path.join(out_dir, "summary-hotspot.md"), "w").write(h)
    open(os.path.join(out_dir, "summary-spread.md"), "w").write(s)
    print(h + "\n" + s + "\nraw: " + os.path.join(out_dir, "raw.csv"))


if __name__ == "__main__":
    main()
