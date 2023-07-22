package interceptor

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/internal"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/server/consumer"
)

func Logging(log logger.Wrapper) consumer.Interceptor {
	return func(ctx context.Context, msg *consumer.Message, handler consumer.Handler) error {
		loggerFields := []logger.Field{
			logger.String(logger.KeyRequestID, msg.MessageId),
			logger.String(logger.KeyRMQExchange, msg.Exchange),
			logger.String(logger.KeyRMQRoutingKey, msg.RoutingKey),
			logger.ByteString(logger.KeyRMQMsgBody, msg.Body),
		}

		if requestID, ok := msg.Headers[internal.HeaderRequestID].(string); ok {
			loggerFields = append(loggerFields, logger.String(logger.KeyRequestID, requestID))
		}

		loggerFields = append(loggerFields, logger.Time(logger.KeyRMQHandlerStartTime, time.Now()))
		return handler(
			ctxzap.ToContext(ctx, log.With(append(loggerFields, ctxzap.TagsToFields(ctx)...)...).GetLogger()),
			msg,
		)
	}
}
