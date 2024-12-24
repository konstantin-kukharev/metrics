package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/konstantin-kukharev/metrics/domain"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"github.com/konstantin-kukharev/metrics/internal/repository/memory"

	"github.com/konstantin-kukharev/metrics/domain/entity"
	ucase "github.com/konstantin-kukharev/metrics/domain/usecase/metric"
)

type TestHandler struct {
	get   *MetricGet
	getV2 *MetricGetV2
	add   *MetricAdd
	addV2 *MetricAddV2
	list  *metricIndex
}

func (h *TestHandler) Get(t, n string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/value/"+domain.MetricGauge+"/test", http.NoBody)
	req.SetPathValue("type", t)
	req.SetPathValue("name", n)

	wr := httptest.NewRecorder()
	h.get.ServeHTTP(wr, req)

	return wr
}

func (h *TestHandler) GetV2(t, n string) *httptest.ResponseRecorder {
	r1, _ := entity.NewMetric(n, t, "")
	data, _ := json.Marshal(r1)
	b := bytes.NewBuffer(data)

	req := httptest.NewRequest(http.MethodPost, "/value", b)
	req.Header.Add("Content-Type", "application/json")

	wr := httptest.NewRecorder()
	h.getV2.ServeHTTP(wr, req)

	return wr
}

func (h *TestHandler) List() *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

	wr := httptest.NewRecorder()
	h.list.ServeHTTP(wr, req)

	return wr
}

func (h *TestHandler) AddV2(t, n, v string) *httptest.ResponseRecorder {
	r1, _ := entity.NewMetric(n, t, v)
	data, _ := json.Marshal(r1)
	b := bytes.NewBuffer(data)

	req := httptest.NewRequest(http.MethodPost, "/update", b)
	req.Header.Add("Content-Type", "application/json")

	wr := httptest.NewRecorder()
	h.addV2.ServeHTTP(wr, req)

	return wr
}

func (h *TestHandler) Add(t, n, v string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/update/"+domain.MetricGauge+"/test/1", http.NoBody)
	req.SetPathValue("type", t)
	req.SetPathValue("name", n)
	req.SetPathValue("val", v)

	wr := httptest.NewRecorder()
	h.add.ServeHTTP(wr, req)

	return wr
}

func NewTestHandler() *TestHandler {
	th := new(TestHandler)
	log := logger.NewSlog()
	store := memory.NewStorage(log)
	g := ucase.NewGetMetric(store)
	a := ucase.NewAddMetric(store)
	l := ucase.NewListMetric(store)
	th.get = NewGetMetric(g)
	th.getV2 = NewMetricGetV2(g)
	th.add = NewAddMetric(a)
	th.addV2 = NewAddMetricV2(a)
	th.list = NewIndexMetric(l)

	return th
}

func TestHandlerGetMetric(t *testing.T) {
	th := NewTestHandler()
	if res := th.Add(domain.MetricCounter, "testCounter", "none"); res.Code != http.StatusBadRequest {
		t.Errorf("got HTTP status code %d, expected 400", res.Code)
	}

	if res := th.Add(domain.MetricGauge, "test", "1"); res.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", res.Code)
	}

	if res := th.Add("incorrect", "test", "1"); res.Code != http.StatusBadRequest {
		t.Errorf("got HTTP status code %d, expected 400", res.Code)
	}

	if res := th.Add("incorrect", "", "1"); res.Code != http.StatusNotFound {
		t.Errorf("got HTTP status code %d, expected 404", res.Code)
	}

	if res := th.Add(domain.MetricGauge, "data", "asdasd"); res.Code != http.StatusBadRequest {
		t.Errorf("got HTTP status code %d, expected 400", res.Code)
	}

	r := th.Get(domain.MetricGauge, "test")
	if r.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", r.Code)
	}

	if !strings.EqualFold(r.Body.String(), "1") {
		t.Errorf(
			`response body "%s" does not equal "1"`,
			r.Body.String(),
		)
	}

	if res := th.Add(domain.MetricGauge, "test", "2"); res.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", res.Code)
	}

	r = th.Get(domain.MetricGauge, "test")
	if r.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", r.Code)
	}
	if !strings.EqualFold(r.Body.String(), "2") {
		t.Errorf(
			`response body "%s" does not equal "2"`,
			r.Body.String(),
		)
	}

	r = th.GetV2(domain.MetricGauge, "test")
	if r.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", r.Code)
	}

	r = th.AddV2(domain.MetricGauge, "test", "3")
	if r.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", r.Code)
	}

	r = th.GetV2(domain.MetricGauge, "test")
	if r.Code != http.StatusOK {
		t.Errorf("got HTTP status code %d, expected 200", r.Code)
	}
}
