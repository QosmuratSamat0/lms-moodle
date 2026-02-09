package main

import (
	"log"

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

func main() {
	cfg := config.Load()

	db := appDeps.InitDatabaseWithConfig(cfg)
	defer db.Close()

	deps := appDeps.BuildDeps(db, cfg)

	modules := appDeps.BuildHTTPModules(deps, cfg.JWTSecret)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	// Global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	for _, module := range modules {
		module.Register(router)
	}

	appInstance := appDeps.New(db)
	defer appInstance.Close()

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
