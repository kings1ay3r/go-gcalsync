package server

import (
	"gcalsync/gophers/clients/logger"
	"gcalsync/gophers/dao"
	"gcalsync/gophers/handlers"
	"gcalsync/gophers/middlewares/auth"
	"gcalsync/gophers/middlewares/response"
	"github.com/gorilla/mux"
	"net/http"
)

func Serve() {
	r := mux.NewRouter()
	dao.InitDB()

	log := logger.NewLogger()

	// Initialize handlers
	handler := handlers.New()

	// Define routes with handlers // TODO: Extract Method, Refactor for Readability/Maintainability
	r.Handle("/connect", auth.BearerAuth(response.APIMiddleware(handler.ConnectHandler))).Methods("GET")
	r.Handle("/callback", auth.BearerAuth(response.APIMiddleware(handler.CallbackHandler))).Methods("GET")
	r.Handle("/index", auth.BearerAuth(response.APIMiddleware(handler.ListEventsHandler))).Methods("GET")

	log.Info(nil, "Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Serve failed to start: %v", err)
	}
}
