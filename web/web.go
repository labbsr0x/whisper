package web

import (
	"log"
	"net/http"
	"time"

	"github.com/abilioesteves/whisper/web/api"
	"github.com/abilioesteves/whisper/web/ui"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Builder defines the parametric information of a whisper server instance
type Builder struct {
	Port       string
	BaseUIPath string
}

// Server holds the information needed to run Whisper
type Server struct {
	*Builder
	UserAPIs api.UserAPI
}

// New builds a Server instance
func (b *Builder) New() (s *Server, err error) {
	s = &Server{}
	s.Builder = b
	s.UserAPIs = new(api.DefaultUserAPI)

	return s, nil
}

// Initialize inits the web server and its apis
func (s *Server) Initialize() error {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/users", s.UserAPIs.ListUsersHandler).Methods("GET")
	router.HandleFunc("/api/users", s.UserAPIs.AddUserHandler).Methods("POST")
	router.HandleFunc("/api/users", s.UserAPIs.RemoveUserHandler).Methods("DELETE")
	router.HandleFunc("/api/users/{clientId}", s.UserAPIs.GetUserHandler).Methods("GET")

	router.Handle("/", ui.Handler(s.BaseUIPath))
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + s.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Infof("Port %v", s.Port)

	log.Print("Initialized")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("server initialization error", err)
		return err
	}
	return nil
}
