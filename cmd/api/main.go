package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ap1-final-mini-moodle/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	appDeps "github.com/ap1-final-mini-moodle/internal/app"
	"github.com/ap1-final-mini-moodle/internal/delivery/http/middleware"
	"github.com/ap1-final-mini-moodle/internal/shared/config"
	"github.com/gin-gonic/gin"
)

// @title Mini Moodle API
// @version 1.0
// @description This is a Mini Moodle API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT token with Bearer prefix. Example: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

func main() {
	cfg := config.Load()

	db := appDeps.InitDatabaseWithConfig(cfg)
	redis := appDeps.InitRedisWithConfig(cfg)
	defer db.Close()
	defer redis.Close()

	deps := appDeps.BuildDeps(db, redis, cfg)

	modules := appDeps.BuildHTTPModules(deps, cfg.JWTSecret)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	for _, module := range modules {
		module.Register(router)
	}

	appInstance := appDeps.New(db)
	defer appInstance.Close()

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutdown signal received, performing graceful shutdown...")

	shCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
