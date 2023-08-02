package statusserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/requiemofthesouls/logger"
	"github.com/requiemofthesouls/monitoring"
	"github.com/requiemofthesouls/postgres"
)

func New(l logger.Wrapper, m monitoring.Wrapper, db postgres.Wrapper, version *version) Manager {
	return &manager{
		l:       l,
		m:       m,
		db:      db,
		version: version,
	}
}

type (
	Manager interface {
		Metrics() http.Handler
		HealthCheck() http.Handler
		Version() http.Handler
	}

	manager struct {
		l       logger.Wrapper
		m       monitoring.Wrapper
		db      postgres.Wrapper
		version *version
	}
)

func (m *manager) Metrics() http.Handler {
	return m.m.MetricsHandler()
}

func (m *manager) HealthCheck() http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if err := m.db.Ping(context.Background()); err != nil {
			m.l.Error("error db.Ping", logger.Error(err))
			http.Error(resp, fmt.Sprintf("db ping error: %v", err.Error()), http.StatusInternalServerError)
			return
		}

		if _, err := resp.Write([]byte("ok")); err != nil {
			m.l.Error("error resp.Write", logger.Error(err))
		}
	})
}
