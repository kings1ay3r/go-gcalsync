package server

import (
	"gcalsync/gophers/dao"
	"gcalsync/gophers/handlers"
	"gcalsync/gophers/middlewares/auth"
	"gcalsync/gophers/middlewares/response"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func Serve() {
	r := mux.NewRouter()
	dao.InitDB()

	// Initialize handlers
	calendarHandler := handlers.NewCalendarHandler()

	// Define routes with handlers
	r.HandleFunc("/connect", handlers.ConnectHandler).Methods("GET")
	r.HandleFunc("/callback", handlers.CallbackHandler).Methods("GET")

	// Define Protected Routes
	r.Handle("/events", auth.BearerAuth(response.ResponseMiddleware(calendarHandler.ListEventsHandler))).Methods("GET")

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Serve failed to start: %v", err)
	}
}
