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
}

func (s *Ping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")

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

func NewPing(cfg pingConfig) *Ping {
	serv := &Ping{dns: cfg.GetDatabaseDNS()}
	return serv
}
