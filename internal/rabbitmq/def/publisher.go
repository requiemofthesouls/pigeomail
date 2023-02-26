package def

import (
	"github.com/requiemofthesouls/pigeomail/internal/rabbitmq"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
	rmqDef "github.com/requiemofthesouls/pigeomail/pkg/modules/rabbitmq/def"
)

const DIAMQPPublisher = "amqp.publisher"

type Publisher = rabbitmq.Publisher

func initAMQPPublisher() container.Def {
	return container.Def{
		Name: DIAMQPPublisher,
		Build: func(cont container.Container) (_ interface{}, err error) {
			var l logDef.Wrapper
			if err = cont.Fill(logDef.DIWrapper, &l); err != nil {
				return nil, err
			}

			var rmq rmqDef.Wrapper
			if err = cont.Fill(rmqDef.DIWrapper, &rmq); err != nil {
				return nil, err
			}

			var p = rabbitmq.NewPublisher(rmq, l)
			return p, nil
		},
	}
}

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(initAMQPPublisher())
	})
}
