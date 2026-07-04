// See wrap_hotspot.js. Runs loadtest/spread.js unchanged.
import def, { options as base } from "../loadtest/spread.js";

export const options = Object.assign({}, base, {
  thresholds: Object.assign({}, base.thresholds, {
    "http_req_duration{expected_response:true}": ["p(95)>=0"],
  }),
  summaryTrendStats: ["avg", "min", "med", "max", "p(90)", "p(95)"],
});

export default def;

export function handleSummary(data) {
  return { [__ENV.SUMMARY_OUT || "summary.json"]: JSON.stringify(data) };
}
