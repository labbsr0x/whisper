package web

import (
	"net/http"
	"time"

	"github.com/labbsr0x/whisper/web/api"

	"github.com/labbsr0x/whisper/web/config"
	"github.com/labbsr0x/whisper/web/ui"

	"github.com/labbsr0x/whisper/web/middleware"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// Server holds the information needed to run Whisper
type Server struct {
	*config.WebBuilder
	UserCredentialsAPIs api.UserCredentialsAPI
	LoginAPIs           api.LoginAPI
	ConsentAPIs         api.ConsentAPI
	HydraAPIs           api.HydraAPI
}

// InitFromWebBuilder builds a Server instance
func (s *Server) InitFromWebBuilder(webBuilder *config.WebBuilder) *Server {
	s.WebBuilder = webBuilder
	s.UserCredentialsAPIs = new(api.DefaultUserCredentialsAPI).InitFromWebBuilder(webBuilder)
	s.LoginAPIs = new(api.DefaultLoginAPI).InitFromWebBuilder(webBuilder)
	s.ConsentAPIs = new(api.DefaultConsentAPI).InitFromWebBuilder(webBuilder)
	s.HydraAPIs = new(api.DefaultHydraAPI).InitFromWebBuilder(webBuilder)

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
	secureRouter := router.PathPrefix("/secure").Subrouter()

	router.PathPrefix("/static").Handler(ui.Handler(s.BaseUIPath)).Methods("GET")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	router.Handle("/login", s.LoginAPIs.LoginGETHandler("/login")).Methods("GET")
	router.Handle("/login", s.LoginAPIs.LoginPOSTHandler()).Methods("POST")

	router.Handle("/consent", s.ConsentAPIs.ConsentGETHandler("/consent")).Methods("GET")
	router.Handle("/consent", s.ConsentAPIs.ConsentPOSTHandler()).Methods("POST")

	router.Handle("/registration", s.UserCredentialsAPIs.GETRegistrationPageHandler("/registration")).Methods("GET")
	router.Handle("/registration", s.UserCredentialsAPIs.POSTHandler()).Methods("POST")

	router.Handle("/email-confirmation", s.UserCredentialsAPIs.GETEmailConfirmationPageHandler("/email-confirmation")).Methods("GET")

	router.Handle("/hydra", s.HydraAPIs.HydraGETHandler("/hydra")).Methods("GET")

	secureRouter.Handle("/update", s.UserCredentialsAPIs.GETUpdatePageHandler("/secure/update")).Methods("GET")
	secureRouter.Handle("/update", s.UserCredentialsAPIs.PUTHandler()).Methods("PUT")

	router.Use(middleware.GetPrometheusMiddleware())
	router.Use(middleware.GetErrorMiddleware())
	secureRouter.Use(s.Self.GetMuxSecurityMiddleware())

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
