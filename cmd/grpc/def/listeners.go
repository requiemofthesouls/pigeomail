package def

import (
	"github.com/requiemofthesouls/container"
	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
	"github.com/requiemofthesouls/pigeomail/internal/grpc/listeners/private_api"
	"github.com/requiemofthesouls/pigeomail/internal/grpc/listeners/public_api"
	"github.com/requiemofthesouls/svc-grpc/server"
	"google.golang.org/grpc"
)

func listenerPublicAPI(cont container.Container) (interface{}, error) {
	return []server.ListenerRegistrant{
		func(srv *grpc.Server) {
			pigeomail_api_pb.RegisterPublicAPIServer(
				srv,
				public_api.New(),
			)
		},
	}, nil
}

func listenerPrivateAPI(cont container.Container) (interface{}, error) {
	return []server.ListenerRegistrant{
		func(srv *grpc.Server) {
			pigeomail_api_pb.RegisterPrivateAPIServer(
				srv,
				private_api.New(),
			)
		},
	}, nil
}
