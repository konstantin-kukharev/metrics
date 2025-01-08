package persistence

import (
	"context"
	"fmt"

	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal/logger"
)

const (
	stateInit    = "init"
	stateRunning = "running"
	statePending = "pending"
	stateStop    = "stop"
)

type MetricStorage struct {
	state string
	log   *logger.Logger

	dns   string
	store *sql.DB
}

func (ms *MetricStorage) Set(ctx context.Context, es ...*entity.Metric) ([]*entity.Metric, error) {
	if ms.state != stateRunning {
		return nil, fmt.Errorf("storage is stopped for new connections")
	}

	tx, err := ms.store.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return nil, err
	}

	for _, e := range es {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO metrics VALUES ($1, $2 ,ROW($3, $4)) ON CONFLICT (name, type) DO UPDATE SET value = metrics.value + EXCLUDED.value;`,
			e.ID, e.MType, e.Delta, e.Value)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	return es, nil
}

func (ms *MetricStorage) connect() error {
	db, err := sql.Open("pgx", ms.dns)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	ms.store = db
	ms.state = stateRunning

	return nil
}

func (ms *MetricStorage) recoverConnection() {
	_ = ms.connect()
}

func (ms *MetricStorage) initialize() error {
	req := `CREATE TYPE mvalue AS (
		delta INTEGER,
		val DOUBLE PRECISION
		);

		CREATE OR REPLACE FUNCTION sum_mvalue(x mvalue, y mvalue) 
		RETURNS mvalue AS $$
		SELECT (y.delta, x.val + y.val);
		$$ language sql;

		CREATE OPERATOR +
		(
		PROCEDURE = sum_mvalue,
		LEFTARG = mvalue,
		RIGHTARG = mvalue
		);

		CREATE TYPE mtype AS ENUM ('gauge','counter');

		CREATE TABLE metrics (
		name VARCHAR NOT NULL,
		type mtype,
		value mvalue
		);

		CREATE unique INDEX metrics_mname_idx ON metrics (name, type);`

	_, err := ms.store.Exec(req)

	return err
}

func (ms *MetricStorage) Run(ctx context.Context) error {
	if err := ms.connect(); err != nil {
		ms.recoverConnection()
	}

	if ms.state != stateRunning {
		return fmt.Errorf("can`t start storage")
	}

	defer ms.store.Close()

	if err := ms.initialize(); err != nil {
		return fmt.Errorf("can`t init storage")
	}

	ms.log.InfoCtx(ctx, "postgres storage is running")

	<-ctx.Done()

	ms.log.InfoCtx(ctx, "postgres storage is stopped")
	ms.state = stateStop
	return nil
}

func NewMetric(l *logger.Logger, dns string) *MetricStorage {
	ms := new(MetricStorage)
	ms.log = l
	ms.state = stateInit
	ms.dns = dns

	return ms
}
