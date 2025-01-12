package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"database/sql"

	"go.uber.org/zap"

	"github.com/jackc/pgconn"
	"github.com/konstantin-kukharev/metrics/domain/entity"
	"github.com/konstantin-kukharev/metrics/internal/logger"
)

const (
	stateInit    = "init"
	stateRunning = "running"
	statePending = "pending"
	stateStop    = "stop"
)

var recoverIntervals = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

type MetricStorage struct {
	state string
	log   *logger.Logger

	dns          string
	store        *sql.DB
	connRecovery chan struct{}
	connErr      chan pgconn.PgError
}

func (ms *MetricStorage) do(ctx context.Context, payload func() error) error {
	if ms.state == statePending {
		select {
		case <-ms.connRecovery:
			break
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	err := payload()
	if err == nil {
		return nil
	}

	// pqErr := err.(pgconn.PgError)
	// if pgerrcode.IsConnectionException(pqErr.Code)
	// 	ms.connErr <- pqErr
	// 	return ms.do(ctx, payload)
	// }

	return err
}

func (ms *MetricStorage) Set(ctx context.Context, es ...*entity.Metric) ([]*entity.Metric, error) {
	err := ms.do(ctx, func() error {
		tx, err := ms.store.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			return err
		}

		sqlInsert := `insert into metrics values ($1, $2 ,$3, $4) on conflict (id, mtype) do update set value = metrics.value + EXCLUDED.value, 
			delta = EXCLUDED.delta`

		for _, e := range es {
			r, err := tx.ExecContext(
				ctx,
				sqlInsert,
				e.ID, e.MType, e.Delta, e.Value)
			if err != nil {
				_ = tx.Rollback()

				return err
			}
			if _, err := r.RowsAffected(); err != nil {
				_ = tx.Rollback()

				return err
			}

			row := tx.QueryRowContext(ctx,
				"select id, mtype, delta, value from metrics where id = $1 AND mtype = $2",
				e.ID, e.MType)
			var d sql.NullInt64
			var i sql.NullFloat64
			err = row.Scan(&e.ID, &e.MType, &d, &i)
			e.SetValue(d, i)
			if err != nil {
				_ = tx.Rollback()

				return err
			}
		}

		err = tx.Commit()
		if err != nil {
			_ = tx.Rollback()

			return err
		}

		return nil
	})

	if err != nil {
		ms.log.ErrorCtx(ctx, "error update", zap.Any("error", err))
		return nil, err
	}

	return es, nil
}

func (ms *MetricStorage) Get(ctx context.Context, ems ...*entity.Metric) ([]*entity.Metric, bool) {
	err := ms.do(ctx, func() error {
		for _, e := range ems {
			var d sql.NullInt64
			var i sql.NullFloat64
			row := ms.store.QueryRowContext(ctx,
				"select id, mtype, delta, value from metrics WHERE id = $1 AND mtype = $2",
				e.ID, e.MType)
			err := row.Scan(&e.ID, &e.MType, &d, &i)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return err
				}

				return fmt.Errorf("can`t get metric: %s", err.Error())
			}

			e.SetValue(d, i)
		}

		return nil
	})

	if err != nil {
		return nil, false
	}

	return ems, true
}

func (ms *MetricStorage) List(ctx context.Context) []*entity.Metric {
	list := make(chan []*entity.Metric, 1)

	err := ms.do(ctx,
		func() error {
			rows, err := ms.store.QueryContext(ctx, "select id, mtype, delta, value from metrics")
			if err != nil {
				close(list)
				return nil
			}

			vals := make([]*entity.Metric, 0)
			for rows.Next() {
				var d sql.NullInt64
				var i sql.NullFloat64
				e := new(entity.Metric)
				err = rows.Scan(&e.ID, &e.MType, &d, &i)
				if err != nil {
					close(list)
					return nil
				}

				e.SetValue(d, i)
				vals = append(vals, e)
			}

			list <- vals
			return nil
		},
	)
	if err != nil {
		return []*entity.Metric{}
	}

	rr, ok := <-list
	if !ok {
		return []*entity.Metric{}
	}

	return rr
}

func (ms *MetricStorage) connect(ctx context.Context) error {
	db, err := sql.Open("pgx", ms.dns)
	if err != nil {
		return err
	}

	ms.store = db
	err = ms.store.PingContext(ctx)
	if err != nil {
		return err
	}

	ms.state = stateRunning

	return nil
}

func (ms *MetricStorage) recoverConnection(ctx context.Context, wait ...time.Duration) error {
	intervals := make(chan struct{})

	go func(ctx context.Context, c chan<- struct{}, w []time.Duration) {
		for _, in := range w {
			select {
			case <-time.After(in):
				c <- struct{}{}
			case <-ctx.Done():
				close(c)
				return
			}
		}
		close(c)
	}(ctx, intervals, wait)

	for {
		select {
		case _, ok := <-intervals:
			if !ok {
				return fmt.Errorf("all db recover attempts are exhausted")
			}
			if err := ms.store.PingContext(ctx); err != nil {
				ms.log.ErrorCtx(ctx, "error while recovering db connection", zap.Any("error", err))
				continue
			}

			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ms *MetricStorage) initialize() error {
	req := `
		DO ' BEGIN
    		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = ''mtype'') THEN
				CREATE TYPE mtype AS ENUM (''gauge'',''counter'');
			END IF;
		END '; 

		CREATE TABLE IF NOT EXISTS metrics (
			id VARCHAR NOT NULL,
			mtype mtype NOT NULL,
			delta INTEGER,
			value DOUBLE PRECISION
		);

		CREATE unique INDEX IF NOT EXISTS metrics_mname_idx ON metrics (id, mtype);`

	_, err := ms.store.Exec(req)

	return err
}

func (ms *MetricStorage) Run(ctx context.Context) error {
	if err := ms.connect(ctx); err != nil {
		return err
	}

	if ms.state != stateRunning {
		return fmt.Errorf("can`t start db storage")
	}

	defer ms.store.Close()

	if err := ms.initialize(); err != nil {
		return fmt.Errorf("can`t init db storage")
	}

	ms.log.InfoCtx(ctx, "postgres storage is running")

	for {
		select {
		case err := <-ms.connErr:
			ms.log.ErrorCtx(ctx, "error db connection", zap.Any("error", err))
			ms.state = statePending
			if err := ms.recoverConnection(ctx, recoverIntervals...); err != nil {
				return err
			}
			for len(ms.connErr) > 0 {
				<-ms.connErr
			}
			ms.state = stateRunning
			close(ms.connRecovery)
			ms.connRecovery = make(chan struct{})
		case <-ctx.Done():
			ms.log.InfoCtx(ctx, "postgres storage is stopped")
			ms.state = stateStop
			return ctx.Err()
		}
	}
}

func NewMetric(l *logger.Logger, dns string) *MetricStorage {
	ms := new(MetricStorage)
	ms.log = l
	ms.state = stateInit
	ms.dns = dns
	ms.connRecovery = make(chan struct{})
	ms.connErr = make(chan pgconn.PgError)

	return ms
}
