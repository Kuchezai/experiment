package app

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"experiment.io/config"
	"experiment.io/internal/controller/http"
	"experiment.io/internal/controller/http/handlers/middleware"
	repo "experiment.io/internal/repo/pg"
	"experiment.io/internal/usecase"
	"experiment.io/pkg/hasher"
	logger "experiment.io/pkg/logger"
	ginLogger "experiment.io/pkg/logger/gin-logger"
	postgres "experiment.io/pkg/storage/pg"
	"github.com/gin-gonic/gin"
)

func Run(ctx context.Context, cfg *config.Config) {
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

	dirToStorageCSV := "./history"
	userRepo, err := repo.NewUserRepository(pg, dirToStorageCSV)
	if err != nil {
		log.Fatal("unable to create user repository")
	}

	// Usecase
	segmentUC := usecase.NewSegmentUsecase(segmentRepo)
	userUC := usecase.NewUserUsecase(userRepo)

	secretKey := cfg.HTTP.JWTSecret
	hasher := hasher.New()
	authUC := usecase.NewAuthUsecase(userRepo, hasher, secretKey)

	// Create http server
	l := logger.New()
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(ginLogger.LoggingMiddleware(l))

	http.SetupRouter(g, l, segmentUC, userUC, authUC, middleware.Authorized(secretKey))
	srv, err := http.NewServer(g, cfg.HTTP)
	if err != nil {
		log.Fatal(err)
	}

	// Start http server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()


	// Graceful shutdown
	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pg.CloseConnections(ctxShutDown); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	<-ctxShutDown.Done()
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
