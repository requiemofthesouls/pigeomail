package def

import (
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
	connDef "github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/connection/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/server"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/server/consumer"
	interceptorDef "github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/server/interceptor/def"
)

const (
	DIManagerPrefix = "rmq.server.manager."
	DIListener      = "rmq.server.listener."
)

type (
	Manager = server.Manager

	Configs []Config
	Config  struct {
		Name                   string `mapstructure:"name"`
		Connection             string `mapstructure:"connection"`
		QueuesHandlersSettings []struct {
			QueueName              string `mapstructure:"queueName"`
			QOS                    uint8  `mapstructure:"qos"`
			NumConsumers           uint8  `mapstructure:"numConsumers"`
			DelayNumFailedAttempts uint8  `mapstructure:"delayNumFailedAttempts"`
			DelayTTLSec            uint16 `mapstructure:"delayTTLSec"`
		} `mapstructure:"queuesHandlersSetting"`
		Restrictions struct {
			AvailableProducts []string `mapstructure:"availableProducts"`
		} `mapstructure:"restrictions"`
	}

	DefinitionBuilder map[string]struct {
		Listener Listener
	}

	Listener func(cont container.Container) (interface{}, error)
)

func (b DefinitionBuilder) AddDefs(builder *container.Builder) error {
	diDefs := make([]container.Def, 0, len(b))
	for name, item := range b {
		name := name
		item := item
		diDefs = append(diDefs, container.Def{
			Name:  DIListener + name,
			Build: item.Listener,
		})
	}

	return builder.Add(diDefs...)
}

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

				interceptors := make([]consumer.Interceptor, 0, len(interceptorDef.List))
				for _, item := range interceptorDef.List {
					var interceptor consumer.Interceptor
					if err := cont.Fill(item, &interceptor); err != nil {
						return nil, err
					}
					interceptors = append(interceptors, interceptor)
				}

				var listenerRegistrants []server.ListenerRegistrant
				if err := cont.Fill(DIListener+cfg.Name, &listenerRegistrants); err != nil {
					return nil, err
				}

				queuesHandlersSettings := make(map[string]*server.QueueHandlerSettings, len(cfg.QueuesHandlersSettings))
				for _, settings := range cfg.QueuesHandlersSettings {
					if _, ok := queuesHandlersSettings[settings.QueueName]; ok {
						continue
					}

					queuesHandlersSettings[settings.QueueName] = &server.QueueHandlerSettings{
						QOS:                    settings.QOS,
						NumConsumers:           settings.NumConsumers,
						DelayNumFailedAttempts: settings.DelayNumFailedAttempts,
						DelayTTL:               int32(settings.DelayTTLSec) * 1000,
					}
				}

				return server.NewManager(
					log.With(logger.String(logger.KeyRMQServerName, cfg.Name)),
					conn,
					listenerRegistrants,
					interceptors,
					queuesHandlersSettings,
					cfg.Restrictions.AvailableProducts,
				)
			},
		})
	}
	return builder.Add(diDefs...)
}
