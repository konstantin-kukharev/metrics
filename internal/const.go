package internal

import "time"

const MetricGauge = "gauge"
const MetricCounter = "counter"

const DefaultServerAddr = "localhost:8080"
const DefaultPoolInterval = time.Duration(2 * time.Second)
const DefaultReportInterval = time.Duration(10 * time.Second)
