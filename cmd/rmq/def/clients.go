package def

import (
	"github.com/requiemofthesouls/container"
	pigeomailpb "github.com/requiemofthesouls/pigeomail/pb"

	clientDef "github.com/requiemofthesouls/svc-rmq/client/def"
)

const (
	DIClientPublisherEvents = "rmq.client.publisher_events"
)

type (
	PublisherEventsClient = pigeomailpb.PublisherEventsRMQClient
)

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name: DIClientPublisherEvents,
				Build: func(cont container.Container) (interface{}, error) {
					var c clientDef.Manager
					if err := cont.Fill(clientDef.DIManagerPrefix+"publisher", &c); err != nil {
						return nil, err
					}

					return pigeomailpb.NewPublisherEventsRMQClient(c), nil
				},
			},
		)
	})
}
