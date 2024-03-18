package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Flgado/fitnessStudioApp/config"
	_ "github.com/Flgado/fitnessStudioApp/docs"
	dbfactory "github.com/Flgado/fitnessStudioApp/internal/database/dbFactory"
	"github.com/Flgado/fitnessStudioApp/routes"
	"github.com/Flgado/fitnessStudioApp/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @tittle FitnessStudioApp
// @version 1
// @Description "App to book"

// @contact.name Joao Folgado
// @contact.url  https://github.com/Flgado
// @contact.email jfolgado94@gmail.com

// @host localhost:8080
func main() {
	configPath := utils.GetConfigPath()

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	portString := cfg.Server.Port

	if portString == "" {
		log.Fatal("PORT is not found in the conf file")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	// database factory
	df := dbfactory.NewDBFactory(cfg)
	dbPoll, err := df.GetDbContext()

	if err != nil {
		log.Fatalf("Impossible to start database pool connections: Error %s", err)
	}

	uRoute := routes.BuildUserRoutes(dbPoll)
	cRoute := routes.BuildClassesRoutes(dbPoll)
	rRoute := routes.BuildReservationRoutes(dbPoll)

	router.Mount("/v1/fitnessstudio/users", uRoute)
	router.Mount("/v1/fitnessstudio/classes", cRoute)
	router.Mount("/v1/fitnessstudio/bookings", rRoute)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Port:", portString)
}
