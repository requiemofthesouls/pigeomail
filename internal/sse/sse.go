package sse

import (
	"encoding/json"
	"net/http"

	"github.com/r3labs/sse/v2"
)

func NewServer() *Server {
	srv := sse.New()
	srv.AutoStream = true
	return &Server{
		server: srv,
	}
}

type (
	Server struct {
		server *sse.Server
	}
)

func (s *Server) Publish(streamID string, object interface{}) error {
	var (
		data []byte
		err  error
	)
	if data, err = json.Marshal(object); err != nil {
		return err
	}
	s.server.Publish(streamID, &sse.Event{Data: data})

	return nil
}

func (s *Server) CreateStream(streamID string) {
	if !s.server.StreamExists(streamID) {
		s.server.CreateStream(streamID)
	}
}

func (s *Server) RemoveStream(streamID string) {
	s.server.RemoveStream(streamID)
}

func (s *Server) StreamExists(streamID string) bool {
	return s.server.StreamExists(streamID)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(w, r)
}
