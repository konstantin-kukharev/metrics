package metric

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/domain/entity"
)

type MetricGetV2 struct {
	service MetricReader
}

func (s *MetricGetV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if err = data.Validate(); err != nil && !errors.Is(err, domain.ErrEmptyMetricValue) {
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

	m, ok := s.service.Do(data)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result := &entity.Metric{
		ID:    m.ID,
		MType: m.MType,
	}
	*result.Delta = *m.Delta
	*result.Value = *m.Value

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

func NewMetricGetV2(srv MetricReader) *MetricGetV2 {
	serv := &MetricGetV2{service: srv}
	return serv
}
