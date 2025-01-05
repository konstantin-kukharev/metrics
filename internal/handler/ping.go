package handler

import (
	"net/http"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type pingConfig interface {
	GetDatabaseDNS() string
}

type Ping struct {
	dns string
	log log
}

type log interface {
	Info(msg string, fields ...any)
}

func (s *Ping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.log.Info("new request", "DB DNS:", s.dns)
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

	w.WriteHeader(http.StatusOK)
}

func NewPing(cfg pingConfig, l log) *Ping {
	serv := &Ping{dns: cfg.GetDatabaseDNS(), log: l}
	return serv
}
