package def

import (
	"net/http"

	"github.com/requiemofthesouls/container"
	logDef "github.com/requiemofthesouls/logger/def"
	monDef "github.com/requiemofthesouls/monitoring/def"
	pgDef "github.com/requiemofthesouls/postgres/def"
	httpServerDef "github.com/requiemofthesouls/svc-http/server/def"

	ss "github.com/requiemofthesouls/pigeomail/internal/http/handlers/status-server"
)

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: httpServerDef.DIHandlerPrefix + "status-server",
			Build: func(cont container.Container) (interface{}, error) {
				var l logDef.Wrapper
				if err := cont.Fill(logDef.DIWrapper, &l); err != nil {
					return nil, err
				}

				var m monDef.Wrapper
				if err := cont.Fill(monDef.DIWrapper, &m); err != nil {
					return nil, err
				}

				var db pgDef.Wrapper
				if err := cont.Fill(pgDef.DIWrapper, &db); err != nil {
					return nil, err
				}

				var (
					statusServer = ss.New(l, m, db, ss.GetVersionFromParams(params))
					mux          = http.NewServeMux()
				)
				mux.Handle("/metrics", statusServer.Metrics())
				mux.Handle("/health", statusServer.HealthCheck())
				mux.Handle("/version", statusServer.Version())

				return mux, nil
			},
		})
	})
}
