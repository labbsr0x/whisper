package web

import (
	"context"
	"net/http"
	"os"
	"os/signal"
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

	router.Handle("/change-password/step-1", s.UserCredentialsAPIs.GETChangePasswordStep1PageHandler("/change-password")).Methods("GET")
	router.Handle("/change-password/step-2", s.UserCredentialsAPIs.GETChangePasswordStep2PageHandler("/change-password")).Methods("GET")
	router.Handle("/change-password", s.UserCredentialsAPIs.POSTChangePasswordPageHandler("/change-password")).Methods("POST")
	router.Handle("/change-password", s.UserCredentialsAPIs.PUTChangePasswordPageHandler("/change-password")).Methods("PUT")

	router.Handle("/hydra", s.HydraAPIs.HydraGETHandler()).Methods("GET")

	secureRouter.Handle("/update", s.UserCredentialsAPIs.GETUpdatePageHandler("/secure/update")).Methods("GET")
	secureRouter.Handle("/update", s.UserCredentialsAPIs.PUTHandler()).Methods("PUT")

	router.Use(middleware.GetPrometheusMiddleware())
	router.Use(middleware.GetErrorMiddleware())
	secureRouter.Use(s.Self.GetMuxSecurityMiddleware())

	return s.ListenAndServe(router)
}

func (s *Server) ListenAndServe(router *mux.Router) error {

	srv := &http.Server{
		Handler:      router,
		Addr:         "0.0.0.0:" + s.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatalf("server initialization error: %v", err)
		}
	}()
	logrus.Info("Server Started")

	<-channel

	logrus.Info("Server Stopped")

	logrus.Debugf("Waiting at most %v seconds for a graceful shutdown", time.Second*s.ShutdownTime)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*s.ShutdownTime)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("server finalization error: %v", err)
	}

	logrus.Info("Server Exited Properly")

	return nil
}
