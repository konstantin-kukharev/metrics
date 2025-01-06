package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type MetricAddV2 struct {
	service metricWriter
}

func (s *MetricAddV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)

	w.Header().Add("Content-Type", "application/json")
	data := &entity.Metric{}
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(data)
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

	if err = data.Validate(); err != nil {
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

	res, err := s.service.Set(r.Context(), data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(res) != 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	m := res[0]
	result := &entity.Metric{
		ID:    m.ID,
		MType: m.MType,
	}

	if data.Delta != nil {
		result.Delta = m.Delta
	}
	if data.Value != nil {
		result.Value = m.Value
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		resp["message"] = "Bad Request. " + err.Error()
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resultJSON)
}

func NewAddMetricV2(srv metricWriter) *MetricAddV2 {
	serv := &MetricAddV2{service: srv}
	return serv
}
