package app

import (
	"fmt"
	"log"
	"net/url"

	"experiment.io/config"
	"experiment.io/internal/controller/http"
	repo "experiment.io/internal/repo/pg"
	"experiment.io/internal/usecase"
	postgres "experiment.io/pkg/storage/pg"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	// PostgreSQL
	pg, err := postgres.New(
		generateDBURL(&cfg.DB, "postgres"),
		postgres.MaxPoolSize(cfg.DB.PoolSize),
		postgres.ConnAttempts(cfg.DB.ConnAttempts),
		postgres.ConnTimeout(cfg.DB.ConnTimeout),
	)
	if err != nil {
		log.Fatal("unable to connect pg")
	}
	// Repository
	segmentRepo := repo.NewSegmentRepository(pg)

	// Usecase
	segmentUC := usecase.NewSegmentUsecase(segmentRepo)

	// Create and start http server
	g := gin.New()
	http.SetupRouter(g, segmentUC)
	srv, err := http.NewServer(g, cfg.HTTP)
	if err != nil {
		log.Fatal(err)
	}
	if err != srv.ListenAndServe() {
		log.Fatal(err)
	}

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
