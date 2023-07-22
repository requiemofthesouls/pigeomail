package interceptor

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/svc-rmq/server/consumer"
)

func Recovery(log logger.Wrapper) consumer.Interceptor {
	return func(ctx context.Context, msg *consumer.Message, handler consumer.Handler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				b, _ := json.Marshal(r)

				log.Error("rmq handle panic",
					logger.String(logger.KeyRequestID, msg.MessageId),
					logger.String(logger.KeyRMQExchange, msg.Exchange),
					logger.String(logger.KeyRMQRoutingKey, msg.RoutingKey),
					logger.ByteString(logger.KeyRMQMsgBody, msg.Body),
					logger.ByteString(logger.KeyRMQHandlerPanicMsg, b),
				)
				err = errors.New("rmq handle panic")
			}
		}()

		return handler(ctx, msg)
	}
}
