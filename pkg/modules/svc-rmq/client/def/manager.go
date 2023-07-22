package def

import (
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"

	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/client"
	connDef "github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/connection/def"
)

const (
	DIManagerPrefix = "rmq.client.manager."
)

type (
	Manager = client.Manager

	Configs []Config
	Config  struct {
		Name       string   `mapstructure:"name"`
		Connection string   `mapstructure:"connection"`
		Routes     []string `mapstructure:"routes"`
	}
)

func (cfgs Configs) AddDefs(builder *container.Builder) error {
	diDefs := make([]container.Def, 0, len(cfgs))
	for _, cfg := range cfgs {
		cfg := cfg
		diDefs = append(diDefs, container.Def{
			Name: DIManagerPrefix + cfg.Name,
			Build: func(cont container.Container) (interface{}, error) {
				var log logDef.Wrapper
				if err := cont.Fill(logDef.DIWrapper, &log); err != nil {
					return nil, err
				}

				var conn connDef.Manager
				if err := cont.Fill(connDef.DIManagerPrefix+cfg.Connection, &conn); err != nil {
					return nil, err
				}

				return client.NewManager(log, conn, cfg.Routes)
			},
		})
	}
	return builder.Add(diDefs...)
}
