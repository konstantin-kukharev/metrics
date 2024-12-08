package domain

import "errors"

var (
	ErrWrongMetricType = errors.New("wrong metric type")
	ErrWrongMetricName = errors.New("wrong metric name")
	ErrInvalidData     = errors.New("invalid data")
)
