package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/SergeiSkv/goose/v3"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
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
	store        *pgx.Conn
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
	query := `insert into metrics (id, mtype, delta, value) values (@id, @mtype ,@delta, @value) ` +
		`on conflict (id, mtype) do update set value = EXCLUDED.value, delta = metrics.delta + EXCLUDED.delta ` +
		`returning id, mtype, delta, value`

	list := make(chan []*entity.Metric, 1)

	err := ms.do(ctx, func() error {
		keysForSelect := make(map[struct {
			k entity.MType
			t string
		}]*entity.Metric)

		batch := &pgx.Batch{}
		for _, e := range es {
			args := pgx.NamedArgs{
				"id": e.ID, "mtype": e.MType,
				"delta": e.Delta, "value": e.Value,
			}
			batch.Queue(query, args)

			keysForSelect[struct {
				k entity.MType
				t string
			}{e.MType, e.ID}] = &entity.Metric{}
		}

		results := ms.store.SendBatch(ctx, batch)
		defer results.Close()

		for c := 0; c < len(es); c++ {
			tmp := &entity.Metric{}
			row := results.QueryRow()
			if err := row.Scan(&tmp.ID, &tmp.MType, &tmp.Delta, &tmp.Value); err != nil {
				return err
			}

			keysForSelect[struct {
				k entity.MType
				t string
			}{tmp.MType, tmp.ID}] = tmp
		}

		ret := make([]*entity.Metric, 0, len(es))
		for _, k := range keysForSelect {
			ret = append(ret, k)
		}

		list <- ret
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
		query := `select id, mtype, delta, value from metrics where id = @id AND mtype = @mtype`
		batch := &pgx.Batch{}
		for _, e := range ems {
			args := pgx.NamedArgs{"id": e.ID, "mtype": e.MType}
			batch.Queue(query, args)
		}

		results := ms.store.SendBatch(ctx, batch)
		defer results.Close()

		for _, e := range ems {
			row := results.QueryRow()
			err := row.Scan(&e.ID, &e.MType, &e.Delta, &e.Value)
			if err != nil {
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
	query := "select id, mtype, delta, value from metrics"

	err := ms.do(ctx,
		func() error {
			rows, err := ms.store.Query(ctx, query)
			if err != nil || rows.Err() != nil {
				close(list)
				return nil
			}
			defer rows.Close()

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
	db, err := pgx.Connect(context.Background(), ms.dns)
	if err != nil {
		return err
	}

	ms.store = db
	err = ms.store.Ping(ctx)
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
			if err := ms.store.Ping(ctx); err != nil {
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

	defer ms.store.Close(ctx)

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
	return &MetricStorage{
		log:          l,
		state:        stateInit,
		dns:          dns,
		connRecovery: make(chan struct{}),
		connErr:      make(chan *pgconn.PgError),
	}
}
