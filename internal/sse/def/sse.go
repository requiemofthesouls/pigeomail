package def

import (
	cfgDef "github.com/requiemofthesouls/config/def"
	"github.com/requiemofthesouls/container"
	"github.com/requiemofthesouls/pigeomail/internal/sse"
)

const (
	DISSEServer = "sse.server"
)

type (
	Server = *sse.Server
)

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name: DISSEServer,
				Build: func(cont container.Container) (interface{}, error) {
					var cfg cfgDef.Wrapper
					if err := cont.Fill(cfgDef.DIWrapper, &cfg); err != nil {
						return nil, err
					}

					return sse.NewServer(), nil
				},
				Close: func(obj interface{}) error {
					return nil
				},
			},
		)
	})
}
