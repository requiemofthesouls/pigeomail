package def

import (
	"errors"

	cfgDef "github.com/requiemofthesouls/config/def"
	"github.com/requiemofthesouls/container"
	logDef "github.com/requiemofthesouls/logger/def"
	repDef "github.com/requiemofthesouls/pigeomail/internal/repository/def"
	"github.com/requiemofthesouls/pigeomail/internal/telegram"
	tgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/telegram/def"
	"github.com/requiemofthesouls/pigeomail/pkg/state"
)

const (
	DITelegramBot = "telegram.bot"
)

type (
	Bot = telegram.Bot
)

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name: DITelegramBot,
				Build: func(cont container.Container) (interface{}, error) {
					var l logDef.Wrapper
					if err := cont.Fill(logDef.DIWrapper, &l); err != nil {
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

					var usersRep repDef.TelegramUsers
					if err := cont.Fill(repDef.DIDBRepositoryTelegramUsers, &usersRep); err != nil {
						return nil, err
					}

					var tgWrapper tgDef.Wrapper
					if err := cont.Fill(tgDef.DITelegramWrapper, &tgWrapper); err != nil {
						return nil, err
					}

					return telegram.NewBot(tgWrapper, l, usersRep, state.NewState(), smtpDomain), nil
				},
				Close: func(obj interface{}) error {
					return nil
				},
			},
		)
	})
}
