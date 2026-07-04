// Runs the versioned loadtest/hotspot.js unchanged, but dumps the full k6
// summary as JSON (SUMMARY_OUT) and forces the 200-only latency submetric.
import def, { options as base } from "../loadtest/hotspot.js";

export const options = Object.assign({}, base, {
  thresholds: Object.assign({}, base.thresholds, {
    "http_req_duration{expected_response:true}": ["p(95)>=0"], // always passes; just forces tracking
  }),
  summaryTrendStats: ["avg", "min", "med", "max", "p(90)", "p(95)"],
});

export default def;

export function handleSummary(data) {
  return { [__ENV.SUMMARY_OUT || "summary.json"]: JSON.stringify(data) };
}
