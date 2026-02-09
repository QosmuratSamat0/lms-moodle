package main

import (
	"log"

	appDeps "github.com/ap1-final-mini-moodle/internal/app"
	"github.com/ap1-final-mini-moodle/internal/shared/config"
	"github.com/gin-gonic/gin"
)

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

	for _, module := range modules {
		module.Register(router)
	}

	appInstance := appDeps.New(db)
	defer appInstance.Close()

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
