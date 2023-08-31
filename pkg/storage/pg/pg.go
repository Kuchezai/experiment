package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	poolSize     int
	connAttempts int
	connTimeout  time.Duration

	*pgxpool.Pool
}

const (
	defaultPoolSize    = 5
	defaultConnAttempt = 3
	defaultConnTimeout = time.Second * 3
)

func New(url string, opts ...Option) (*Postgres, error) {
	op := "storage.pg.New"

	pg := &Postgres{
		poolSize:     defaultPoolSize,
		connAttempts: defaultConnAttempt,
		connTimeout:  defaultConnTimeout,
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", op, err)
	}

	for _, opt := range opts {
		opt(pg)
	}

	for ; pg.connAttempts > 0; pg.connAttempts-- {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			break
		}

		time.Sleep(pg.connTimeout)
	}

	return pg, nil
}

func (p *Postgres) CloseConnections(ctx context.Context) error {
	op := "storage.pg.CloseConnections"
	
	done := make(chan struct{})
	go func() {
		p.Close()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%s, %w", op, errors.New(("time expired")))
	}
}
