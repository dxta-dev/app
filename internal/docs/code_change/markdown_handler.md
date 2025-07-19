---
title: "Code Change Metrics"
---
The Code Change engineering metric quantifies the team’s weekly development activity by measuring the total number of lines of code added, modified, or deleted across our repositories.

* **Source**: Computed from commit diffs in our Git version-control system, excluding merge commits and auto-generated files.
* **Aggregation**: Grouped by ISO week (Monday–Sunday).
* **Purpose**:
   * Tracks engineering velocity and throughput over time.
   * Highlights spikes (e.g., major feature work or refactors) and troughs (e.g., stabilization periods, planning, or holidays).
   * Helps correlate process changes (code freezes, new tooling) with fluctuations in developer output.
