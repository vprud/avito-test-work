package app

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/cut4cut/avito-test-work/config"
	v1 "github.com/cut4cut/avito-test-work/internal/controller/http/v1"
	"github.com/cut4cut/avito-test-work/internal/usecase"
	"github.com/cut4cut/avito-test-work/internal/usecase/repo"
	"github.com/cut4cut/avito-test-work/pkg/logger"
	"github.com/cut4cut/avito-test-work/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	r := repo.New(pg)
	accountUseCase := usecase.New(r)

	// HTTP Server
	handler := gin.Default()
	v1.NewRouter(handler, l, *accountUseCase)

	handler.Run()
}
