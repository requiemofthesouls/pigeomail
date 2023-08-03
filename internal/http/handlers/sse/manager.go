package sse

import (
	"net/http"

	"github.com/requiemofthesouls/logger"
	sseDef "github.com/requiemofthesouls/pigeomail/internal/sse/def"
)

func New(
	l logger.Wrapper,
	sseServer sseDef.Server,
) Manager {
	return &manager{
		l:         l,
		sseServer: sseServer,
	}
}

type (
	Manager interface {
		Stream() http.Handler
	}

	manager struct {
		l         logger.Wrapper
		sseServer sseDef.Server
	}
)

func (m *manager) Stream() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		streamID := r.URL.Query().Get("stream")
		if streamID == "" {
			http.Error(w, "Please specify a stream!", http.StatusInternalServerError)
			return
		}

		go func() {
			// Received Browser Disconnection
			<-r.Context().Done()
			m.l.Info("The client is disconnected here")
			// Remove Stream
			m.sseServer.RemoveStream(streamID)
			return
		}()

		m.l.Info("new connection")

		m.sseServer.ServeHTTP(w, r)
	})
}
