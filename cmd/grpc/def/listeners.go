package def

import (
	"errors"

	cfgDef "github.com/requiemofthesouls/config/def"
	"github.com/requiemofthesouls/container"
	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
	"github.com/requiemofthesouls/pigeomail/internal/grpc/listeners/private_api"
	"github.com/requiemofthesouls/pigeomail/internal/grpc/listeners/public_api"
	sseDef "github.com/requiemofthesouls/pigeomail/internal/sse/def"
	"github.com/requiemofthesouls/svc-grpc/server"
	"google.golang.org/grpc"
)

func listenerPublicAPI(cont container.Container) (interface{}, error) {
	var sse sseDef.Server
	if err := cont.Fill(sseDef.DISSEServer, &sse); err != nil {
		return nil, err
	}

	var cfg cfgDef.Wrapper
	if err := cont.Fill(cfgDef.DIWrapper, &cfg); err != nil {
		return nil, err
	}

	var smtpDomain string
	if smtpDomain = cfg.GetString("smtp_domain"); smtpDomain == "" {
		return nil, errors.New("smtp_domain is empty")
	}

	return []server.ListenerRegistrant{
		func(srv *grpc.Server) {
			pigeomail_api_pb.RegisterPublicAPIServer(
				srv,
				public_api.New(sse, smtpDomain),
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
