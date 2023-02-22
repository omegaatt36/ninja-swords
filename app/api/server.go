package api

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/omegaatt36/ninja-swords/appmodule/bot"
	"github.com/omegaatt36/ninja-swords/logging"
	ginprometheus "github.com/omegaatt36/ninja-swords/pkg/go-gin-prometheus"
)

// Server is a HTTP server.
type Server struct {
}

// RegisterMiddleware registers middleware for all endpoints.
func (s *Server) RegisterMiddleware(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	config.AllowAllOrigins = false

	r.Use(cors.New(config))
	r.Use(RateLimit())
}

// RegisterEndpoint installs api representation layer processing function.
func (s *Server) RegisterEndpoint(r *gin.Engine) {
	r.POST("echo", func(ctx *gin.Context) {
		var req bot.EchoRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx.JSON(200, req)
	})
}

// Start starts HTTP server.
func (s *Server) Start(ctx context.Context, apiAddr string) {
	gin.ForceConsoleColor()

	// setup gin.
	apiEngine := gin.New()
	apiEngine.RedirectTrailingSlash = true

	prome := ginprometheus.NewPrometheus("api")
	prome.Use(apiEngine)

	s.RegisterMiddleware(apiEngine)

	// setup endpoint.
	s.RegisterEndpoint(apiEngine)

	srv := &http.Server{
		Addr:    apiAddr,
		Handler: apiEngine,
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown: ", err)
		}
	}()

	logging.Get().Info("starts serving...")
	if err := srv.ListenAndServe(); err != nil &&
		!errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
