package http

import (
    "net/http"
    "github.com/gorilla/mux"
    "your_project/internal/engine" // Adjust the import path as necessary
)

type Server struct {
    engine *engine.Engine
}

func NewServer(engine *engine.Engine) *Server {
    return &Server{engine: engine}
}

func (s *Server) Start(addr string) error {
    router := mux.NewRouter()
    s.routes(router)

    return http.ListenAndServe(addr, router)
}

func (s *Server) routes(router *mux.Router) {
    // Define your routes here
    router.HandleFunc("/api/v1/tasks", s.createTaskHandler).Methods("POST")
    router.HandleFunc("/api/v1/tasks/{id}", s.getTaskHandler).Methods("GET")
    router.HandleFunc("/api/v1/tasks/{id}/retry", s.retryTaskHandler).Methods("POST")
    // Add more routes as needed
}

func (s *Server) createTaskHandler(w http.ResponseWriter, r *http.Request) {
    // Handler logic for creating a task
}

func (s *Server) getTaskHandler(w http.ResponseWriter, r *http.Request) {
    // Handler logic for retrieving a task
}

func (s *Server) retryTaskHandler(w http.ResponseWriter, r *http.Request) {
    // Handler logic for retrying a task
}