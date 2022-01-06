package server

import (
	"net/http"
)

func (s *Server) getAlivenessHandler() http.HandlerFunc {
	def := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	if s.alivenessCheck != nil {
		return s.alivenessCheck(def)
	}
	return def
}

func (s *Server) getReadinessHandler() http.HandlerFunc {
	def := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	if s.readinessCheck != nil {
		return s.readinessCheck(def)
	}
	return def
}

func (s *Server) getHealthCheck() http.HandlerFunc {
	def := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK!"))
	})

	if s.healthCheck != nil {
		return s.healthCheck(def)
	}
	return def
}
