package def

import (
	"errors"

	"github.com/requiemofthesouls/pigeomail/pkg/modules/config"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
)

const DIWrapper = "config.wrapper"

type (
	Wrapper = config.Wrapper
)

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		var ok bool
		if _, ok = params["config"]; !ok {
			return errors.New("can't get required parameter config path")
		}

		var path string
		if path, ok = params["config"].(string); !ok {
			return errors.New(`parameter "config_path" should be string`)
		}

		return builder.Add(container.Def{
			Name: DIWrapper,
			Build: func(container container.Container) (_ interface{}, err error) {
				return config.New(path)
			},
		})
	})
}
