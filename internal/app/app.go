package app

import (
	"fmt"
	"log"
	nethttp "net/http"

	"github.com/cezarychodun/freellms/internal/config"
	"github.com/cezarychodun/freellms/internal/database"
	apphttp "github.com/cezarychodun/freellms/internal/http"
	"github.com/cezarychodun/freellms/internal/modules/usage"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// App has router and db instances.
type App struct {
	Router *mux.Router
	DB     *sqlx.DB
}

// Initialize initializes the app with predefined configuration.
func (a *App) Initialize(config *config.Config) {
	db, err := database.Open(config.DB)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	modelResourcesRepository := usage.NewModelResourcesRepository(db)

	a.DB = db
	a.Router = apphttp.NewRouter(modelResourcesRepository)
	fmt.Println("App initialized successfully")
}

// Run starts the HTTP server.
func (a *App) Run(host string) {
	defer a.DB.Close()

	log.Fatal(nethttp.ListenAndServe(host, a.Router))
}
