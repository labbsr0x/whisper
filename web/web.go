package web

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/abilioesteves/careless-whisper/web/api"
	"github.com/abilioesteves/careless-whisper/web/ui"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Initialize inits the web server and its apis
func Initialize() error {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/users", api.ListUsersHandler).Methods("GET")
	router.HandleFunc("/api/users", api.AddUserHandler).Methods("POST")
	router.HandleFunc("/api/users", api.RemoveUserHandler).Methods("DELETE")
	router.HandleFunc("/api/users/{clientId}", api.GetUserHandler).Methods("GET")

	router.Handle("/login", ui.Handler())
	router.Handle("/consent", ui.Handler())
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + os.Getenv("PORT"),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Print("Initialized")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("server initializing error", err)
		return err
	}
	return nil
}
