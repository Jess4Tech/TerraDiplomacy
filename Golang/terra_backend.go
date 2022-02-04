package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"jess.buetow/terra_backend/authorization"
	"jess.buetow/terra_backend/config"
	"jess.buetow/terra_backend/database"
)

func main() {
	config.LoadConfiguration()

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	authRouter := apiRouter.PathPrefix("/auth").Subrouter()

	authManager := authorization.NewAuthorizationManager()

	authRouter.HandleFunc("/status", authManager.StatusHttp).Methods(http.MethodGet)

	authRouter.HandleFunc("/login", authManager.LoginHttp).Methods(http.MethodPost)

	authRouter.HandleFunc("/logout", authManager.LogoutHttp).Methods(http.MethodPost)

	authRouter.HandleFunc("/otac", func(w http.ResponseWriter, r *http.Request) {
		if authManager.Authenticated(r, authorization.Server) {
			authManager.Otac(w, r)
		} else if authManager.Present(r) {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}).Methods(http.MethodPost)

	databaseHandler, err := database.NewDatabaseHandler()
	if err != nil {
		log.Fatalln(err)
	}

	apiRouter.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "":
		case "GET":
			if authManager.Authenticated(r, authorization.Player) {
				// log.Println("Getting projects")
				databaseHandler.GetProjectsHttp(w, r)
			} else if authManager.Present(r) {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		case "POST":
			if authManager.Authenticated(r, authorization.Admin) {
				// log.Printf("Adding project")
				databaseHandler.AddProjectHttp(w, r)
			} else if authManager.Present(r) {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		case "DELETE":
			if authManager.Authenticated(r, authorization.Admin) {
				// log.Printf("Deleting project")
				databaseHandler.DeleteProjectHttp(w, r)
			} else if authManager.Present(r) {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		}
	}).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)

	apiRouter.HandleFunc("/tension", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "":
		case "GET":
			if authManager.Authenticated(r, authorization.Player) {
				// log.Println("Getting tension")
				databaseHandler.LeaderboardTensionHttp(w, r)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		case "POST":
			if authManager.Authenticated(r, authorization.Server) {
				// log.Println("Adding tension")
				databaseHandler.SetTensionHttp(w, r)
			} else if authManager.Present(r) {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		case "DELETE":
			if authManager.Authenticated(r, authorization.Server) {
				// log.Println("Deleting tension")
				databaseHandler.DeleteTensionHttp(w, r)
			} else if authManager.Present(r) {
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		}
	}).Methods(http.MethodGet, http.MethodPost, http.MethodDelete)

	finished := make(chan os.Signal, 1)
	signal.Notify(finished, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logRouter := handlers.CombinedLoggingHandler(os.Stdout, router)
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{config.Config.FrontendAddress.SchemedAddress()},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
		MaxAge:           60 * 60 * 24 * 7,
		Debug:            true,
	})

	server := &http.Server{
		Addr:    config.Config.ListenAddress.Address(),
		Handler: handlers.RecoveryHandler()(cors.Handler(logRouter)),
	}
	go func() {
		if config.Config.ListenAddress.Secure() {
			if err := server.ListenAndServeTLS(config.Config.SslConfig.Certificate, config.Config.SslConfig.Key); err != nil && err != http.ErrServerClosed {
				log.Fatalln(err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalln(err)
			}
		}
	}()
	log.Printf("Listening on %s", config.Config.ListenAddress.SchemedAddress())

	<-finished
	log.Println("Shutting down server (This may take up to five seconds)")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server shutdown successfully")
}
