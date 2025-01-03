package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type MetricAddV2 struct {
	service MetricWriter
}

func (s *MetricAddV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	resp := make(map[string]string)
	if headerContentTtype != "application/json" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		resp["message"] = "Content Type is not application/json"
		jsonResp, _ := json.Marshal(resp)
		_, _ = w.Write(jsonResp)
		return
	}

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

	if err = s.service.Do(data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	result := &entity.Metric{
		ID:    data.ID,
		MType: data.MType,
	}

	if data.Delta != nil {
		result.Delta = data.Delta
	}
	if data.Value != nil {
		result.Value = data.Value
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

func NewAddMetricV2(srv MetricWriter) *MetricAddV2 {
	serv := &MetricAddV2{service: srv}
	return serv
}
