package web

import (
	"log"
	"net/http"
	"time"

	"github.com/abilioesteves/whisper/misc"

	"github.com/abilioesteves/whisper/web/middleware"

	"github.com/abilioesteves/whisper/web/api"
	"github.com/abilioesteves/whisper/web/ui"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Builder defines the parametric information of a whisper server instance
type Builder struct {
	Port          string
	BaseUIPath    string
	HydraEndpoint string
}

// Server holds the information needed to run Whisper
type Server struct {
	*Builder
	UserAPIs    api.UserAPI
	LoginAPIs   api.LoginAPI
	ConsentAPIs api.ConsentAPI
}

// New builds a Server instance
func (b *Builder) New() (s *Server, err error) {
	s = &Server{}
	hydraClient := new(misc.HydraClient).Init(b.HydraEndpoint)

	s.Builder = b
	s.UserAPIs = new(api.DefaultUserAPI)
	s.LoginAPIs = new(api.DefaultLoginAPI).Init(hydraClient, b.BaseUIPath)
	s.ConsentAPIs = new(api.DefaultConsentAPI).Init(hydraClient, b.BaseUIPath)

	return s, nil
}

// Initialize inits the web server and its apis
func (s *Server) Initialize() error {
	router := mux.NewRouter().StrictSlash(true)
	secureRouter := router.PathPrefix("/api").Subrouter()

	router.PathPrefix("/static").Handler(ui.Handler(s.BaseUIPath)).Methods("GET")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")
	router.Handle("/login", http.StripPrefix("/login", s.LoginAPIs.LoginGETHandler())).Methods("GET")
	router.HandleFunc("/login", s.LoginAPIs.LoginPOSTHandler).Methods("POST")
	router.Handle("/consent", http.StripPrefix("/consent", s.ConsentAPIs.ConsentGETHandler())).Methods("GET")
	router.HandleFunc("/consent", s.ConsentAPIs.ConsentPOSTHandler).Methods("POST")

	secureRouter.HandleFunc("/users", s.UserAPIs.ListUsersHandler).Methods("GET")
	secureRouter.HandleFunc("/users", s.UserAPIs.AddUserHandler).Methods("POST")
	secureRouter.HandleFunc("/users", s.UserAPIs.RemoveUserHandler).Methods("DELETE")
	secureRouter.HandleFunc("/users/{clientId}", s.UserAPIs.GetUserHandler).Methods("GET")

	router.Use(middleware.PrometheusMiddleware)
	router.Use(middleware.ErrorMiddleware)
	secureRouter.Use(middleware.SecurityMiddleware)

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
