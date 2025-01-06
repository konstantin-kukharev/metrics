package handler

import (
	"net/http"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/konstantin-kukharev/metrics/internal/logger"
	"go.uber.org/zap"
)

type Ping struct {
	dns string
	log *logger.Logger
}

func (s *Ping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.log.InfoCtx(r.Context(), "new request:",
		zap.String("DB DNS", s.dns),
	)
	if s.dns == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("pgx", s.dns)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func NewPing(dns string, l *logger.Logger) *Ping {
	serv := &Ping{dns: dns, log: l}
	return serv
}
