package def

import (
	cfgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/config/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
)

const DIWrapper = "logger.wrapper"

type Wrapper = logger.Wrapper

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DIWrapper,
			Build: func(container container.Container) (_ interface{}, err error) {
				var cfg cfgDef.Wrapper
				if err = container.Fill(cfgDef.DIWrapper, &cfg); err != nil {
					return nil, err
				}

				var loggerCfg logger.Config
				if err = cfg.UnmarshalKey("logger", &loggerCfg); err != nil {
					return nil, err
				}

				return logger.New(
					loggerCfg,
					[]logger.Field{
						logger.String("namespace", cfg.GetString("namespace")),
						logger.String("service", cfg.GetString("service")),
						logger.String("instance", cfg.GetString("instance")),
					},
				)
			},
			Close: func(obj interface{}) error {
				_ = obj.(logger.Wrapper).Sync()
				return nil
			},
		})
	})
}
