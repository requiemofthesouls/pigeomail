package def

import (
	"context"

	cfgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/config/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/rabbitmq"
)

const (
	DIWrapper = "rabbitmq.wrapper"
)

type (
	Wrapper = rabbitmq.Wrapper
)

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name: DIWrapper,
				Build: func(container container.Container) (interface{}, error) {
					var cfg cfgDef.Wrapper
					if err := container.Fill(cfgDef.DIWrapper, &cfg); err != nil {
						return nil, err
					}

					var rCfg rabbitmq.Config
					if err := cfg.UnmarshalKey("rabbitmq", &rCfg); err != nil {
						return nil, err
					}

					var l logDef.Wrapper
					if err := container.Fill(logDef.DIWrapper, &l); err != nil {
						return nil, err
					}

					return rabbitmq.New(context.Background(), rCfg, l)
				},
				Close: func(obj interface{}) error {
					if err := obj.(rabbitmq.Wrapper).CloseChannel(); err != nil {
						return err
					}

					return obj.(rabbitmq.Wrapper).CloseConnection()
				},
			})
	})
}
