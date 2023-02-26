package def

import (
	"context"

	cfgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/config/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/mongo"
)

const (
	DIWrapper = "mongo.wrapper"
)

type (
	Wrapper = mongo.Wrapper
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

					var mCfg mongo.Config
					if err := cfg.UnmarshalKey("mongo", &mCfg); err != nil {
						return nil, err
					}

					return mongo.New(context.Background(), mCfg)
				},
				Close: func(obj interface{}) error {
					return obj.(mongo.Wrapper).Close()
				},
			})
	})
}
