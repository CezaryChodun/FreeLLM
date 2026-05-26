package app

import (
	"fmt"
	"log"
	nethttp "net/http"

	"github.com/cezarychodun/freellms/internal/config"
	"github.com/cezarychodun/freellms/internal/database"
	apphttp "github.com/cezarychodun/freellms/internal/http"
	"github.com/cezarychodun/freellms/internal/modules/modelgroups"
	"github.com/cezarychodun/freellms/internal/modules/models"
	"github.com/cezarychodun/freellms/internal/modules/ratelimits"
	"github.com/cezarychodun/freellms/internal/modules/usage"
	"github.com/cezarychodun/freellms/internal/proxy"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type App struct {
	Router *mux.Router
	DB     *sqlx.DB
}

func (a *App) Initialize(config *config.Config) {
	db, err := database.Open(config.DB)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	modelRepo := models.NewModelRepository(db)
	rateLimitRepo := ratelimits.NewRateLimitRepository(db)
	usageRepo := usage.NewModelResourcesRepository(db)
	modelGroupRepo := modelgroups.NewModelGroupRepository(db)

	if err := ratelimits.LoadConfig(modelRepo, rateLimitRepo, modelGroupRepo, "config.yml", "defaults"); err != nil {
		log.Fatalf("failed to load rate limits config: %v", err)
	}

	selector := proxy.NewModelSelector(modelRepo, rateLimitRepo, usageRepo, modelGroupRepo)

	a.DB = db
	a.Router = apphttp.NewRouter(usageRepo, selector, modelGroupRepo)
	fmt.Println("App initialized successfully")
}

func (a *App) Run(host string) {
	defer a.DB.Close()
	log.Fatal(nethttp.ListenAndServe(host, a.Router))
}
