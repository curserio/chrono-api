package server

import "github.com/curserio/chrono-api/pkg/logger"

const (
	defaultLanguage = "en"
)

func WithLogger(l logger.Logger) func(*Server) {
	return func(s *Server) {
		s.logger = l
	}
}

func WithDefaultLanguage(lang string) func(*Server) {
	return func(s *Server) {
		s.defaultLanguage = lang
	}
}
