package internal

import "time"

const MetricGauge = "gauge"
const MetricCounter = "counter"
const CacheKey = "PollCount"

const DefaultServerAddr = "localhost:8080"
const DefaultPoolInterval = 2 * time.Second
const DefaultReportInterval = 10 * time.Second
