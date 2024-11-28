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

func Test_server_MetricUpdate(t *testing.T) {
	type want struct {
		code        int
		contentType string
	}
	type params struct {
		Type  string
		Name  string
		Value string
	}
	tests := []struct {
		name    string
		link    string
		pathVal params
		want    want
	}{
		{
			name:    "add gauge",
			link:    "/update/gauge/test1/1234/",
			pathVal: params{Type: internal.MetricGauge, Name: "test", Value: "1234"},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "add counter",
			link:    "/update/counter/test1/1234/",
			pathVal: params{Type: internal.MetricCounter, Name: "test", Value: "1234"},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain",
			},
		},
		{
			name:    "no name",
			link:    "/update/gauge/test1/1234/",
			pathVal: params{Type: internal.MetricGauge, Name: "", Value: "1234"},
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:    "no type",
			link:    "/update/gauge/test1/1234/",
			pathVal: params{Type: "", Name: "test", Value: "1234"},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
		{
			name:    "no data",
			link:    "/update/gauge/test1/1234/",
			pathVal: params{Type: "", Name: "", Value: ""},
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain",
			},
		},
		{
			name:    "invalid data",
			link:    "/update/gauge/test1/1234/",
			pathVal: params{Type: internal.MetricGauge, Name: "test1", Value: "asdsad"},
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemStorage()

			serv := service.NewMetric(store, dto.Gauge(), dto.Counter())
			srv := NewMetric(serv)

			request := httptest.NewRequest(http.MethodPost, tt.link, nil)
			request.SetPathValue("type", tt.pathVal.Type)
			request.SetPathValue("name", tt.pathVal.Name)
			request.SetPathValue("val", tt.pathVal.Value)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			srv.MetricUpdate(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
