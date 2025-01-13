package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"database/sql"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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
	connErr      chan *pgconn.PgError
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

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgerrcode.IsConnectionException(pgErr.Code) {
		ms.connErr <- pgErr
		return ms.do(ctx, payload)
	}

	return err
}

func (ms *MetricStorage) Set(ctx context.Context, es ...*entity.Metric) ([]*entity.Metric, error) {
	list := make(chan []*entity.Metric, 1)

	err := ms.do(ctx, func() error {
		tx, err := ms.store.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})
		if err != nil {
			return err
		}

		sqlInsert, err := tx.PrepareContext(ctx, `insert into metrics values ($1, $2 ,$3, $4) `+
			`on conflict (id, mtype) do update set value = EXCLUDED.value, delta = metrics.delta + EXCLUDED.delta `+
			`returning id, mtype, delta, value`)
		if err != nil {
			return err
		}
		defer sqlInsert.Close()
		keysForSelect := make(map[struct{ k, t string }]*entity.Metric)

		for _, e := range es {
			row := sqlInsert.QueryRowContext(
				ctx,
				e.ID, e.MType, e.Delta, e.Value)

			ne := new(entity.Metric)
			err := row.Scan(&ne.ID, &ne.MType, &ne.Delta, &ne.Value)
			if err != nil {
				_ = tx.Rollback()

				return err
			}

			keysForSelect[struct {
				k string
				t string
			}{ne.ID, ne.MType}] = ne
		}

		err = tx.Commit()
		if err != nil {
			_ = tx.Rollback()

			return err
		}

		results := make([]*entity.Metric, 0, len(keysForSelect))
		for _, value := range keysForSelect {
			results = append(results, value)
		}

		list <- results
		return nil
	})

	if err != nil {
		ms.log.ErrorCtx(ctx, "error update", zap.Any("error", err))
		close(list)

		return nil, err
	}

	return <-list, nil
}

func (ms *MetricStorage) Get(ctx context.Context, ems ...*entity.Metric) ([]*entity.Metric, bool) {
	err := ms.do(ctx, func() error {
		sqlGet, err := ms.store.PrepareContext(ctx, `select id, mtype, delta, value from metrics where id = $1 AND mtype = $2`)
		if err != nil {
			return err
		}
		defer sqlGet.Close()
		for _, e := range ems {
			row := sqlGet.QueryRowContext(ctx, e.ID, e.MType)
			err := row.Scan(&e.ID, &e.MType, &e.Delta, &e.Value)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return err
				}

				return fmt.Errorf("can`t get metric: %s", err.Error())
			}
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
			if err != nil || rows.Err() != nil {
				close(list)
				return nil
			}

			vals := make([]*entity.Metric, 0)
			for rows.Next() {
				e := new(entity.Metric)
				err = rows.Scan(&e.ID, &e.MType, &e.Delta, &e.Value)
				if err != nil {
					close(list)
					return nil
				}
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

// TODO: [change to goose migrations](https://github.com/pressly/goose)
func (ms *MetricStorage) initialize() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(ms.store, migrationsDir); err != nil {
		return err
	}

	return nil
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
	ms.connErr = make(chan *pgconn.PgError)

	return ms
}
