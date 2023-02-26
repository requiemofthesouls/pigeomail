package def

import (
	"github.com/requiemofthesouls/pigeomail/internal/rabbitmq"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
	rmqDef "github.com/requiemofthesouls/pigeomail/pkg/modules/rabbitmq/def"
)

const DIAMQPConsumer = "amqp.consumer"

type Consumer = rabbitmq.Consumer

func initAMQPConsumer() container.Def {
	return container.Def{
		Name: DIAMQPConsumer,
		Build: func(cont container.Container) (_ interface{}, err error) {
			var l logDef.Wrapper
			if err = cont.Fill(logDef.DIWrapper, &l); err != nil {
				return nil, err
			}

			var rmq rmqDef.Wrapper
			if err = cont.Fill(rmqDef.DIWrapper, &rmq); err != nil {
				return nil, err
			}

			return rabbitmq.NewConsumer(rmq, l), nil
		},
	}
}

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(initAMQPConsumer())
	})
}
