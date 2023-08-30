package app

import (
	"fmt"
	"log"
	"net/url"

	"experiment.io/config"
	"experiment.io/internal/controller/http"
	repo "experiment.io/internal/repo/pg"
	"experiment.io/internal/usecase"
	logger "experiment.io/pkg/logger"
	ginLogger "experiment.io/pkg/logger/gin-logger"
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
	dirToStorageCSV := "./history"
	// Repository
	segmentRepo := repo.NewSegmentRepository(pg)
	userRepo, err := repo.NewUserRepository(pg, dirToStorageCSV)
	if err != nil {
		log.Fatal("unable to create user repository")
	}

	// Usecase
	segmentUC := usecase.NewSegmentUsecase(segmentRepo)
	userUC := usecase.NewUserUsecase(userRepo)

	// Create and start http server
	l := logger.New()
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(ginLogger.LoggingMiddleware(l))
	http.SetupRouter(g, l, segmentUC, userUC)
	srv, err := http.NewServer(g, cfg.HTTP)
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.ListenAndServe(); err != nil {
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
