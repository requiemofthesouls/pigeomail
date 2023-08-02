package def

import (
	"github.com/requiemofthesouls/container"
	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
	"github.com/requiemofthesouls/svc-grpc/server"
	serverDef "github.com/requiemofthesouls/svc-grpc/server/def"
)

var grpcServerDefs = map[string]serverDef.DefinitionBuilder{
	"public-api": {
		Listener: listenerPublicAPI,
		Gateways: serverDef.Gateways{
			RegistrantList: []server.GatewayRegistrant{pigeomail_api_pb.RegisterPublicAPIHandler},
			HttpStatusMap:  server.HTTPStatusMap{},
		},
		Middlewares: serverDef.Middlewares{
			Unary:  serverDef.PublicClientMiddlewares,
			Stream: serverDef.PublicClientMiddlewares,
		},
	},
	"private-api": {
		Listener: listenerPrivateAPI,
		Gateways: serverDef.Gateways{
			RegistrantList: []server.GatewayRegistrant{},
			HttpStatusMap:  server.HTTPStatusMap{},
		},
		Middlewares: serverDef.Middlewares{
			Unary:  serverDef.DefaultMiddlewares,
			Stream: serverDef.DefaultMiddlewares,
		},
	},
}

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return serverDef.AddDefinitions(builder, grpcServerDefs)
	})
}
