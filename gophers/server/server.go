package server

import (
	"fmt"
	"gcalsync/gophers/clients/logger"
	"gcalsync/gophers/dao"
	"gcalsync/gophers/handlers"
	"gcalsync/gophers/middlewares/auth"
	"gcalsync/gophers/middlewares/response"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func Serve() {
	r := mux.NewRouter()
	err := dao.InitDB()
	if err != nil {
		logger.GetInstance().Error(nil, "unable to init database: %v", err)
		return
	}

	log := logger.NewLogger()

	// Initialize handlers
	handler, err := handlers.New()
	if err != nil {
		logger.GetInstance().Error(nil, "unable to init handler: %v", err)
		return
	}

	// Define routes with handlers // TODO: Extract Method, Refactor for Readability/Maintainability
	r.Handle("/connect", auth.BearerAuth(response.APIMiddleware(handler.ConnectHandler))).Methods("GET")
	r.Handle("/callback", auth.BearerAuth(response.APIMiddleware(handler.CallbackHandler))).Methods("GET")
	r.Handle("/index", auth.BearerAuth(response.APIMiddleware(handler.ListEventsHandler))).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	portString := fmt.Sprintf(":%s", port)
	log.Info(nil, "Starting server on :%s...", portString)
	if err := http.ListenAndServe(portString, r); err != nil {
		log.Fatalf("Serve failed to start: %v", err)
	}

}
