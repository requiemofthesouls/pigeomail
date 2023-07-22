package def

import (
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/server/interceptor"
)

const (
	DIInterceptorRecovery = "rmq.server.middleware.recovery"
	DIInterceptorLogging  = "rmq.server.middleware.logging"
)

var List = []string{
	DIInterceptorLogging,
	DIInterceptorRecovery,
}

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name: DIInterceptorRecovery,
				Build: func(cont container.Container) (interface{}, error) {
					var log logDef.Wrapper
					if err := cont.Fill(logDef.DIWrapper, &log); err != nil {
						return nil, err
					}

					return interceptor.Recovery(log), nil
				},
			},
			container.Def{
				Name: DIInterceptorLogging,
				Build: func(cont container.Container) (interface{}, error) {
					var log logDef.Wrapper
					if err := cont.Fill(logDef.DIWrapper, &log); err != nil {
						return nil, err
					}

					return interceptor.Logging(log), nil
				},
			},
		)
	})
}
