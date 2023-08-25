package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	poolSize     int
	connAttempts int
	connTimeout  time.Duration

	Pool *pgxpool.Pool
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
