package web

import (
	"net/http"
	"time"

	"github.com/abilioesteves/whisper/web/api"

	"github.com/abilioesteves/whisper/web/config"
	"github.com/abilioesteves/whisper/web/ui"

	"github.com/abilioesteves/whisper/web/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Server holds the information needed to run Whisper
type Server struct {
	*config.WebBuilder
	UserAPIs    api.UserAPI
	LoginAPIs   api.LoginAPI
	ConsentAPIs api.ConsentAPI
}

// InitFromWebBuilder builds a Server instance
func (s *Server) InitFromWebBuilder(webBuilder *config.WebBuilder) *Server {
	s.WebBuilder = webBuilder
	s.UserAPIs = new(api.DefaultUserAPI)
	s.LoginAPIs = new(api.DefaultLoginAPI).InitFromWebBuilder(webBuilder)
	s.ConsentAPIs = new(api.DefaultConsentAPI).InitFromWebBuilder(webBuilder)

	logLevel, err := logrus.ParseLevel(s.LogLevel)
	if err != nil {
		logrus.Errorf("Not able to parse log level string. Setting default level: info.")
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	return s
}

// Run initializes the web server and its apis
func (s *Server) Run() error {
	router := mux.NewRouter().StrictSlash(true)
	secureRouter := router.PathPrefix("/api").Subrouter()

	router.PathPrefix("/static").Handler(ui.Handler(s.BaseUIPath)).Methods("GET")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	router.Handle("/login", s.LoginAPIs.LoginGETHandler("/login")).Methods("GET")
	router.Handle("/login", s.LoginAPIs.LoginPOSTHandler()).Methods("POST")

	router.Handle("/consent", s.ConsentAPIs.ConsentGETHandler("/consent")).Methods("GET")
	router.Handle("/consent", s.ConsentAPIs.ConsentPOSTHandler()).Methods("POST")

	secureRouter.HandleFunc("/users", s.UserAPIs.AddUserHandler).Methods("POST")
	secureRouter.HandleFunc("/users", s.UserAPIs.RemoveUserHandler).Methods("DELETE")
	secureRouter.HandleFunc("/users/{userID}", s.UserAPIs.UpdateUserHandler).Methods("PUT")

	router.Use(middleware.PrometheusMiddleware)
	router.Use(middleware.ErrorMiddleware)
	secureRouter.Use(middleware.SecurityMiddleware)

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + s.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Info("Initialized")
	err := srv.ListenAndServe()
	if err != nil {
		logrus.Fatal("server initialization error", err)
		return err
	}
	return nil
}
