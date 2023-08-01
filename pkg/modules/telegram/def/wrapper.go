package def

import (
	cfgDef "github.com/requiemofthesouls/config/def"
	"github.com/requiemofthesouls/container"
	logDef "github.com/requiemofthesouls/logger/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/telegram"
)

const (
	DITelegramWrapper = "telegram.wrapper"
)

type (
	Wrapper = telegram.Wrapper
)

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name: DITelegramWrapper,
				Build: func(cont container.Container) (interface{}, error) {
					var cfg cfgDef.Wrapper
					if err := cont.Fill(cfgDef.DIWrapper, &cfg); err != nil {
						return nil, err
					}

					var l logDef.Wrapper
					if err := cont.Fill(logDef.DIWrapper, &l); err != nil {
						return nil, err
					}

					var tgConfig telegram.Config
					if err := cfg.UnmarshalKey("telegram", &tgConfig); err != nil {
						return nil, err
					}

					return telegram.New(&tgConfig, l)
				},
				Close: func(obj interface{}) error {
					return nil
				},
			},
		)
	})
}
