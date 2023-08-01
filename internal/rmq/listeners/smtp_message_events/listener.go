package smtp_message_events

import (
	"context"

	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
)

func NewListener(
	smtpMessageEventHandlerV1 *SMTPMessageEventHandlerV1,
) *listener {
	return &listener{
		smtpMessageEventHandlerV1: smtpMessageEventHandlerV1,
	}
}

type listener struct {
	smtpMessageEventHandlerV1 *SMTPMessageEventHandlerV1
}

func (l *listener) SMTPMessageV1(ctx context.Context, ev *pigeomail_api_pb.SMTPMessageEventV1) error {
	return l.smtpMessageEventHandlerV1.Handle(ctx, ev)
}
