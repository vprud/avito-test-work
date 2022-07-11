// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// Swagger docs.
	_ "github.com/cut4cut/avito-test-work/docs"
	"github.com/cut4cut/avito-test-work/internal/usecase"
	"github.com/cut4cut/avito-test-work/pkg/logger"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// NewRouter -.
// Swagger spec:
// @title       Avito test work API
// @description Simple transaction service
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, u usecase.AccountUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")

	h1 := handler.Group("/swagger", CORSMiddleware())
	{
		h1.GET("/*any", swaggerHandler)
	}

	h2 := handler.Group("/v1")
	{
		newAccountRoutes(h2, u, l)
	}
}
