package main

import (
	"log"
	"os"

	appDeps "github.com/ap1-final-mini-moodle/internal/app"
	"github.com/ap1-final-mini-moodle/internal/shared/config"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := appDeps.InitDatabaseWithConfig(cfg)
	defer db.Close()

	deps := appDeps.BuildDeps(db)

	modules := appDeps.BuildHTTPModules(deps, cfg.JWTSecret)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	for _, module := range modules {
		module.Register(router)
	}

	appInstance := appDeps.New(db)
	appInstance.startBackgroundWorkers()
	defer appInstance.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s (env: %s)\n", port, cfg.Env)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
