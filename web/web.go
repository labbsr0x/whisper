package web

import (
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
	Port               string
	BaseUIPath         string
	HydraAdminEndpoint string
	LogLevel           string
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
	hydraClient := new(misc.HydraClient).Init(b.HydraAdminEndpoint)

	s.Builder = b
	s.UserAPIs = new(api.DefaultUserAPI)
	s.LoginAPIs = new(api.DefaultLoginAPI).Init(hydraClient, b.BaseUIPath)
	s.ConsentAPIs = new(api.DefaultConsentAPI).Init(hydraClient, b.BaseUIPath)

	logLevel, err := logrus.ParseLevel(s.LogLevel)
	if err != nil {
		logrus.Errorf("Not able to parse log level string. Setting Default.")
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	return s, nil
}

// Initialize inits the web server and its apis
func (s *Server) Initialize() error {
	router := mux.NewRouter().StrictSlash(true)
	secureRouter := router.PathPrefix("/api").Subrouter()

	router.PathPrefix("/static").Handler(ui.Handler(s.BaseUIPath)).Methods("GET")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	router.Handle("/login", s.LoginAPIs.LoginGETHandler("/login")).Methods("GET")
	router.Handle("/login", s.LoginAPIs.LoginPOSTHandler()).Methods("POST")

	router.Handle("/consent", s.ConsentAPIs.ConsentGETHandler("/consent")).Methods("GET")
	router.Handle("/consent", s.ConsentAPIs.ConsentPOSTHandler()).Methods("POST")

	router.HandleFunc("/users", s.UserAPIs.AddUserHandler).Methods("POST")

	secureRouter.HandleFunc("/users", s.UserAPIs.ListUsersHandler).Methods("GET")
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

	logrus.Info("Initialized")
	err := srv.ListenAndServe()
	if err != nil {
		logrus.Fatal("server initialization error", err)
		return err
	}
	return nil
}
