package app

import (
	"fmt"
	"net/url"

	"experiment.io/config"
	"experiment.io/pkg/storage/pg"
)

func Run(cfg *config.Config) {
	storage, err := pg.New(
		generateDBURL(&cfg.DB, "postgres"),
		pg.MaxPoolSize(cfg.DB.PoolSize),
		pg.ConnAttempts(cfg.DB.ConnAttempts),
		pg.ConnTimeout(cfg.DB.ConnTimeout),
	)

	fmt.Println(storage, err)
}

func generateDBURL(config *config.DB, scheme string) string {

	u := url.URL{
		Scheme: scheme,
		User:   url.UserPassword(config.User, config.Pass),
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		Path:   config.Name,
	}

	return u.String()
}
