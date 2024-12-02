package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	dto "github.com/konstantin-kukharev/metrics/cmd/server/metric"
	"github.com/konstantin-kukharev/metrics/cmd/server/service"
	"github.com/konstantin-kukharev/metrics/cmd/server/storage"
	"github.com/konstantin-kukharev/metrics/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_server_MetricAdd(t *testing.T) {
	type want struct {
		code        int
		contentType string
		value       string
	}
	type params struct {
		Type  string
		Name  string
		Value string
	}

	store := storage.NewMemStorage()
	serv := service.NewMetric(store, dto.Gauge(), dto.Counter())

	tests := []struct {
		name    string
		pathVal params
		srv     internal.MetricService
		want    want
	}{
		{
			name:    "add gauge",
			pathVal: params{Type: internal.MetricGauge, Name: "test", Value: "1234"},
			srv:     serv,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
				value:       "1234.000",
			},
		},
		{
			name:    "add counter",
			pathVal: params{Type: internal.MetricCounter, Name: "test", Value: "10"},
			srv:     serv,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
				value:       "10",
			},
		},
		{
			name:    "add counter compare",
			pathVal: params{Type: internal.MetricCounter, Name: "test", Value: "10"},
			srv:     serv,
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
				value:       "20",
			},
		},
		{
			name:    "no name",
			pathVal: params{Type: internal.MetricGauge, Name: "", Value: "1234"},
			srv:     serv,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:    "no type",
			pathVal: params{Type: "", Name: "test", Value: "1234"},
			srv:     serv,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name:    "no data",
			pathVal: params{Type: "", Name: "", Value: ""},
			srv:     serv,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:    "invalid data",
			pathVal: params{Type: internal.MetricGauge, Name: "test1", Value: "asdsad"},
			srv:     serv,
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hAdd := NewAddMetric(tt.srv)

			mType := tt.pathVal.Type
			mName := tt.pathVal.Name
			mValue := tt.pathVal.Value

			linkAdd := "/update/" + mType + "/" + mName + "/" + mValue
			requestAdd := httptest.NewRequest(http.MethodPost, linkAdd, http.NoBody)
			requestAdd.SetPathValue("type", mType)
			requestAdd.SetPathValue("name", mName)
			requestAdd.SetPathValue("val", mValue)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			hAdd.ServeHTTP(w, requestAdd)
			res := w.Result()
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			if tt.want.code != http.StatusOK {
				return
			}

			hGet := NewGetMetric(tt.srv)
			linkGet := "/value/" + mType + "/" + mName
			requestGet := httptest.NewRequest(http.MethodGet, linkGet, nil)
			requestGet.SetPathValue("type", mType)
			requestGet.SetPathValue("name", mName)
			// создаём новый Recorder
			w = httptest.NewRecorder()
			hGet.ServeHTTP(w, requestGet)
			resGet := w.Result()
			defer resGet.Body.Close()
			if tt.want.value != "" {
				b, err := io.ReadAll(resGet.Body)
				assert.Nil(t, err)
				assert.Equal(t, tt.want.value, string(b))
			}
		})
	}
}

func Test_server_MetricGet(t *testing.T) {
	type want struct {
		code        int
		contentType string
		value       string
	}
	type params struct {
		Type  string
		Name  string
		Value string
	}

	store := storage.NewMemStorage()
	serv := service.NewMetric(store, dto.Gauge(), dto.Counter())

	tests := []struct {
		name    string
		pathVal params
		srv     internal.MetricService
		want    want
	}{
		{
			name:    "get gauge",
			pathVal: params{Type: internal.MetricGauge, Name: "", Value: "1234"},
			srv:     serv,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:    "get undefined",
			pathVal: params{Type: "", Name: "", Value: "1234"},
			srv:     serv,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:    "get wrong",
			pathVal: params{Type: "wrongtype", Name: "", Value: "1234"},
			srv:     serv,
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mType := tt.pathVal.Type
			mName := tt.pathVal.Name

			hGet := NewGetMetric(tt.srv)
			linkGet := "/value/" + mType + "/" + mName
			requestGet := httptest.NewRequest(http.MethodGet, linkGet, nil)
			requestGet.SetPathValue("type", mType)
			requestGet.SetPathValue("name", mName)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			hGet.ServeHTTP(w, requestGet)
			resGet := w.Result()
			defer resGet.Body.Close()
			if tt.want.value != "" {
				b, err := io.ReadAll(resGet.Body)
				assert.Nil(t, err)
				assert.Equal(t, tt.want.value, string(b))
			}
		})
	}

	newStore := storage.NewMemStorage()
	newServ := service.NewMetric(newStore, dto.Gauge(), dto.Counter())
	linkBrokenGet := "/update/gauge/test/123"
	requestBrokenGet := httptest.NewRequest(http.MethodGet, linkBrokenGet, nil)
	w := httptest.NewRecorder()
	hNewGet := NewAddMetric(newServ)
	hNewGet.ServeHTTP(w, requestBrokenGet)
	resNewGet := w.Result()
	defer resNewGet.Body.Close()
	assert.Equal(t, http.StatusNotFound, resNewGet.StatusCode)
}

func Test_server_MetricList(t *testing.T) {
	type want struct {
		code int
	}

	store := storage.NewMemStorage()
	serv := service.NewMetric(store, dto.Gauge(), dto.Counter())

	tests := []struct {
		name string
		srv  internal.MetricService
		want want
	}{
		{
			name: "list test",
			srv:  serv,
			want: want{
				code: http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hIndex := NewIndexMetric(tt.srv)
			linkGet := "/"
			requestGet := httptest.NewRequest(http.MethodGet, linkGet, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			hIndex.ServeHTTP(w, requestGet)
			resGet := w.Result()
			defer resGet.Body.Close()
			assert.Equal(t, tt.want.code, resGet.StatusCode)
		})
	}
}
