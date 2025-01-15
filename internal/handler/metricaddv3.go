package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type MetricAddV3 struct {
	service metricWriter
}

func (s *MetricAddV3) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	w.Header().Add("Content-Type", "application/json")
	data := make([]*entity.Metric, 0)
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if errors.As(err, &unmarshalErr) {
			resp["message"] = "Bad Request. Wrong Type provided for field " + unmarshalErr.Field
		} else {
			resp["message"] = "Bad Request. " + err.Error()
		}
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)

		return
	}

	for _, m := range data {
		err = m.Validate()
		if err == nil {
			continue
		}

		switch {
		case errors.Is(err, domain.ErrWrongMetricName):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}

		resp["message"] = "Bad Request. " + err.Error()
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)

		return
	}

	res, err := s.service.Set(r.Context(), data...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultJSON, err := json.Marshal(res)
	if err != nil {
		resp["message"] = "Bad Request. " + err.Error()
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resultJSON)
}

func NewAddMetricV3(srv metricWriter) *MetricAddV3 {
	serv := &MetricAddV3{service: srv}
	return serv
}
