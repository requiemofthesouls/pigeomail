package def

import (
	"github.com/requiemofthesouls/container"
	serverDef "github.com/requiemofthesouls/svc-rmq/server/def"
)

var serverDefs = serverDef.DefinitionBuilder{
	"smtp-message-events": {
		Listener: smtpMessageEventsListener,
	},
}

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return serverDefs.AddDefs(builder)
	})
}
