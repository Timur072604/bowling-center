package web

import (
	"Main/bowling"
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	state *bowling.State
}

func New(state *bowling.State) *Server {
	return &Server{state: state}
}

func (s *Server) Run() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/status.html")
	})

	http.HandleFunc("/api/status", s.handleStatusJSON)

	fmt.Println("http://localhost:8080")
	return http.ListenAndServe(":8080", nil)
}

func (s *Server) handleStatusJSON(w http.ResponseWriter, _ *http.Request) {
	snapshot := s.state.GetSnapshot()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(snapshot)
	if err != nil {
		return
	}
}
