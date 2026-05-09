package app

import (
	"database/sql"
	"fmt"
	"log"
	nethttp "net/http"

	"github.com/cezarychodun/freellms/internal/config"
	"github.com/cezarychodun/freellms/internal/database"
	apphttp "github.com/cezarychodun/freellms/internal/http"
	"github.com/gorilla/mux"
)

// App has router and db instances.
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize initializes the app with predefined configuration.
func (a *App) Initialize(config *config.Config) {
	db, err := database.Open(config.DB)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	a.DB = db
	a.Router = apphttp.NewRouter()
	fmt.Println("App initialized successfully")
}

// Run starts the HTTP server.
func (a *App) Run(host string) {
	defer a.DB.Close()

	log.Fatal(nethttp.ListenAndServe(host, a.Router))
}
