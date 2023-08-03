package public_api

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/requiemofthesouls/logger"
	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
	sseDef "github.com/requiemofthesouls/pigeomail/internal/sse/def"
)

func New(
	sse sseDef.Server,
	smtpDomain string,
) pigeomail_api_pb.PublicAPIServer {
	return &manager{
		smtpDomain: smtpDomain,
		clients: clients{
			sse: sse,
		},
	}
}

type (
	manager struct {
		smtpDomain string
		clients    clients
	}

	clients struct {
		sse sseDef.Server
	}
)

func getLogger(ctx context.Context) logger.Wrapper {
	return logger.NewFromZap(ctxzap.Extract(ctx))
}
