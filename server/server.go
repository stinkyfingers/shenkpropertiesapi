package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/stinkyfingers/shenkpropertiesapi/email"
	"github.com/stinkyfingers/shenkpropertiesapi/storage"
)

type Server struct {
	Storage storage.Storage
}

type Response struct {
	Message string `json:"message"`
}

func NewServer(profile string) (*Server, error) {
	store, err := storage.NewS3(profile)
	if err != nil {
		return nil, err
	}

	return &Server{
		Storage: store,
	}, nil
}

// NewMux returns the router
func NewMux(s *Server) (http.Handler, error) {
	mux := http.NewServeMux()
	mux.Handle("/data", cors(s.data))
	mux.Handle("/sendEmail", cors(sendEmail))
	mux.Handle("/images", cors(s.getImages))
	mux.Handle("/file", cors(s.getFile))
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

func (s *Server) data(w http.ResponseWriter, r *http.Request) {
	data, err := s.Storage.Get(storage.DATA_BUCKET, "shenkproperties.json")
	if err != nil {
		errorResponse(w, err)
		return
	}
	var b bytes.Buffer
	_, err = b.ReadFrom(data)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b.Bytes())
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

func (s *Server) getImages(w http.ResponseWriter, r *http.Request) {
	property := r.URL.Query().Get("property")
	keys, err := s.Storage.List(storage.DATA_BUCKET, property)
	if err != nil {
		errorResponse(w, err)
		return
	}
	for i, key := range keys {
		keys[i] = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", storage.DATA_BUCKET, key)
	}
	Photos := struct {
		Keys []string `json:"keys"`
	}{
		keys,
	}
	j, err := json.Marshal(Photos)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (s *Server) getFile(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	suffix := strings.TrimPrefix(filepath.Ext(key), ".")
	data, err := s.Storage.Get(storage.DATA_BUCKET, key)
	if err != nil {
		errorResponse(w, err)
		return
	}
	var b bytes.Buffer
	_, err = b.ReadFrom(data)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.Header().Set("Content-Type", fmt.Sprintf("application/%s", suffix))
	w.Write(b.Bytes())
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
