package server

import (
	"encoding/json"
	"net/http"

	"github.com/stinkyfingers/shenkpropertiesapi/email"
)

type Server struct {
}

func NewServer(profile string) (*Server, error) {
	return &Server{}, nil
}

// NewMux returns the router
func NewMux(s *Server) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.Handle("/sendEmail", cors(sendEmail))
	mux.Handle("/test", cors(status))
	return mux, nil
}

func isPermittedOrigin(origin string) string {
	var permittedOrigins = []string{
		"http://localhost:3000",
		"https://shenkproperties.com",
		"https://www.shenkproperties.com",
		"http://localhost:3001",
	}
	for _, permittedOrigin := range permittedOrigins {
		if permittedOrigin == origin {
			return origin
		}
	}
	return ""
}

func cors(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		permittedOrigin := isPermittedOrigin(r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", permittedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next := http.HandlerFunc(handler)
		next.ServeHTTP(w, r)
	})
}

func status(w http.ResponseWriter, r *http.Request) {
	status := struct {
		Health string `json:"health"`
	}{
		"healthy",
	}
	j, err := json.Marshal(status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(j)
}

func sendEmail(w http.ResponseWriter, r *http.Request) {
	var app email.Application
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = email.SendEmail(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Email Sent"))
}
