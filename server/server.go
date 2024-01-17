package server

import (
	"encoding/json"
	"net/http"

	"github.com/stinkyfingers/shenkpropertiesapi/email"
)

type Server struct {
}

type Response struct {
	Message string `json:"message"`
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
		handler.ServeHTTP(w, r)
	})
}

func status(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Health string `json:"health"`
	}{
		"healthy",
	}
	j, err := json.Marshal(resp)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func sendEmail(w http.ResponseWriter, r *http.Request) {
	var app email.Application
	err := json.NewDecoder(r.Body).Decode(&app)
	if err != nil {
		errorResponse(w, err)
		return
	}
	err = email.SendEmail(app)
	if err != nil {
		errorResponse(w, err)
		return
	}
	j, err := json.Marshal(Response{Message: "Email Sent"})
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func errorResponse(w http.ResponseWriter, err error) {
	errStruct := struct {
		Error string `json:"error"`
	}{
		err.Error(),
	}
	j, err := json.Marshal(errStruct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Error(w, string(j), http.StatusInternalServerError)
}
