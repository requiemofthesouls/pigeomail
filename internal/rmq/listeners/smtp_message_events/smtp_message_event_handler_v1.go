package smtp_message_events

import (
	"context"
	"fmt"
	"html"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/requiemofthesouls/logger"
	pigeomail_api_pb "github.com/requiemofthesouls/pigeomail/api/pb"
	"github.com/requiemofthesouls/pigeomail/internal/repository"
	"github.com/requiemofthesouls/pigeomail/internal/repository/entity"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/telegram"
	"go.uber.org/zap"
)

func NewSMTPMessageEventHandlerV1(
	users repository.TelegramUsers,
	telegram telegram.Wrapper,
) *SMTPMessageEventHandlerV1 {
	return &SMTPMessageEventHandlerV1{
		repositories: smtpMessageEventHandlerV1Repositories{
			users: users,
		},
		clients: smtpMessageEventHandlerV1Clients{
			telegram: telegram,
		},
	}
}

type (
	SMTPMessageEventHandlerV1 struct {
		repositories smtpMessageEventHandlerV1Repositories
		clients      smtpMessageEventHandlerV1Clients
	}

	smtpMessageEventHandlerV1Repositories struct {
		users repository.TelegramUsers
	}

	smtpMessageEventHandlerV1Clients struct {
		telegram telegram.Wrapper
	}
)

func (h *SMTPMessageEventHandlerV1) Handle(ctx context.Context, ev *pigeomail_api_pb.SMTPMessageEventV1) error {
	l := getLogger(ctx)

	l.Debug("smtp message received", zap.String("to", ev.GetTo()))

	var (
		user *entity.TelegramUser
		err  error
	)
	if user, err = h.repositories.users.GetByEMail(ctx, ev.GetTo()); err != nil {
		l.Error("repositories.users.GetByEMail error", zap.Error(err))
		return err
	}

	if !user.IsExist() {
		l.Info("user not found")
		return nil
	}

	textTemplate := `
<b>From:</b> %s
<b>To:</b> %s
<b>Subject:</b> %s
----------------
%s
----------------
`
	if len(ev.GetBody()) > 4096 {
		text := fmt.Sprintf(
			textTemplate,
			html.EscapeString(ev.GetFrom()),
			html.EscapeString(ev.GetTo()),
			html.EscapeString(ev.GetSubject()),
			html.EscapeString(ev.GetBody()[:3000]),
		)

		tgMsg := tgbotapi.NewMessage(user.ChatID, text)
		tgMsg.ParseMode = tgbotapi.ModeHTML

		if _, err = h.clients.telegram.Send(tgMsg); err != nil {
			l.Error("error send message", zap.Error(err))
		}

		for i := 3000; i < len(ev.GetBody()); i += 4096 {
			y := i + 4096
			if y > len(ev.GetBody()) {
				y = len(ev.GetBody())
			}

			tgMsg = tgbotapi.NewMessage(user.ChatID, html.EscapeString(ev.GetBody()[i:y]))
			tgMsg.ParseMode = tgbotapi.ModeHTML

			if _, err = h.clients.telegram.Send(tgMsg); err != nil {
				l.Error("error send message", zap.Error(err))
			}
		}

		return nil
	}

	text := fmt.Sprintf(
		textTemplate,
		html.EscapeString(ev.GetFrom()),
		html.EscapeString(ev.GetTo()),
		html.EscapeString(ev.GetSubject()),
		html.EscapeString(ev.GetBody()),
	)

	tgMsg := tgbotapi.NewMessage(user.ChatID, text)
	tgMsg.ParseMode = tgbotapi.ModeHTML

	if _, err = h.clients.telegram.Send(tgMsg); err != nil {
		l.Error("error send message", zap.Error(err))
	}

	l.Debug("message successfully forwarded to telegram")

	return nil
}

func getLogger(ctx context.Context) logger.Wrapper {
	return logger.NewFromZap(ctxzap.Extract(ctx))
}
