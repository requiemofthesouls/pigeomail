package def

import (
	"github.com/requiemofthesouls/container"
	repDef "github.com/requiemofthesouls/pigeomail/internal/repository/def"
	"github.com/requiemofthesouls/pigeomail/internal/rmq/listeners/smtp_message_events"
	pigeomailpb "github.com/requiemofthesouls/pigeomail/pb"
	tgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/telegram/def"
	"github.com/requiemofthesouls/svc-rmq/server"
)

func smtpMessageEventsListener(cont container.Container) (interface{}, error) {
	var users repDef.TelegramUsers
	if err := cont.Fill(repDef.DIDBRepositoryTelegramUsers, &users); err != nil {
		return nil, err
	}

	var telegram tgDef.Wrapper
	if err := cont.Fill(tgDef.DITelegramWrapper, &telegram); err != nil {
		return nil, err
	}

	return []server.ListenerRegistrant{
		func(srv server.Manager) {
			pigeomailpb.RegisterSMTPMessageEventsRMQServer(
				srv,
				smtp_message_events.NewListener(
					smtp_message_events.NewSMTPMessageEventHandlerV1(
						users,
						telegram,
					),
				),
			)
		},
	}, nil
}
