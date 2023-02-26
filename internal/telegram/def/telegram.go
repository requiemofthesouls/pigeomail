package def

import (
	rmqDef "github.com/requiemofthesouls/pigeomail/internal/rabbitmq/def"
	"github.com/requiemofthesouls/pigeomail/internal/telegram"

	repDef "github.com/requiemofthesouls/pigeomail/internal/repository/def"
	cfgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/config/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
)

const DITGBot = "tg.bot"

type TGBot = telegram.Bot

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name:  DITGBot,
				Build: initTGBot,
			})
	})
}

func initTGBot(cont container.Container) (_ interface{}, err error) {
	var cfg cfgDef.Wrapper
	if err = cont.Fill(cfgDef.DIWrapper, &cfg); err != nil {
		return nil, err
	}

	var l logDef.Wrapper
	if err = cont.Fill(logDef.DIWrapper, &l); err != nil {
		return nil, err
	}

	var repo repDef.EmailState
	if err = cont.Fill(repDef.DIDBRepositoryEmailState, &repo); err != nil {
		return nil, err
	}

	var consumer rmqDef.Consumer
	if err = cont.Fill(rmqDef.DIAMQPConsumer, &consumer); err != nil {
		return nil, err
	}

	var tgConfig telegram.Config
	if err = cfg.UnmarshalKey("telegram", &tgConfig); err != nil {
		return nil, err
	}

	return telegram.NewBot(tgConfig, l, repo, consumer)
}
