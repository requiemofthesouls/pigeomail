package def

import (
	"time"

	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/connection"
)

const (
	DIManagerPrefix = "rmq.connection.manager."
)

type (
	Configs []Config
	Config  struct {
		Name     string       `mapstructure:"name"`
		Host     string       `mapstructure:"host"`
		Port     int32        `mapstructure:"port"`
		Username string       `mapstructure:"username"`
		Password string       `mapstructure:"password"`
		Params   ConfigParams `mapstructure:"params"`
	}
	ConfigParams struct {
		ConnectionName string `mapstructure:"connectionName"`
		HeartbeatSec   uint8  `mapstructure:"heartbeatSec"`
		Locale         string `mapstructure:"locale"`
	}

	Manager = connection.Manager
)

func (cfgs Configs) AddDefs(builder *container.Builder) error {
	diDefs := make([]container.Def, 0, len(cfgs))
	for _, cfg := range cfgs {
		diDefs = append(diDefs, container.Def{
			Name: DIManagerPrefix + cfg.Name,
			Build: func(cont container.Container) (interface{}, error) {
				var log logDef.Wrapper
				if err := cont.Fill(logDef.DIWrapper, &log); err != nil {
					return nil, err
				}

				return connection.NewManager(
					log.With(logger.String(logger.KeyRMQConnectionName, cfg.Name)),
					connection.Config{
						Host:     cfg.Host,
						Port:     cfg.Port,
						Username: cfg.Username,
						Password: cfg.Password,
						Params: connection.ConfigParams{
							ConnectionName: cfg.Params.ConnectionName,
							Heartbeat:      time.Duration(cfg.Params.HeartbeatSec) * time.Second,
							Locale:         cfg.Params.Locale,
						},
					},
				)
			},
		})
	}
	return builder.Add(diDefs...)
}
