package def

import (
	cfgDef "github.com/requiemofthesouls/config/def"
	"github.com/requiemofthesouls/container"
	logDef "github.com/requiemofthesouls/logger/def"
	rmqDef "github.com/requiemofthesouls/pigeomail/cmd/rmq/def"
	"github.com/requiemofthesouls/pigeomail/internal/receiver"
	repDef "github.com/requiemofthesouls/pigeomail/internal/repository/def"
)

const DISMTPReceiver = "smtp.receiver"

type Receiver = receiver.Receiver

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name:  DISMTPReceiver,
				Build: initSMTPReceiver,
			})
	})
}

func initSMTPReceiver(cont container.Container) (_ interface{}, err error) {
	var cfg cfgDef.Wrapper
	if err = cont.Fill(cfgDef.DIWrapper, &cfg); err != nil {
		return nil, err
	}

	var emailRep repDef.TelegramUsers
	if err = cont.Fill(repDef.DIDBRepositoryTelegramUsers, &emailRep); err != nil {
		return nil, err
	}

	var l logDef.Wrapper
	if err = cont.Fill(logDef.DIWrapper, &l); err != nil {
		return nil, err
	}

	var publisher rmqDef.PublisherEventsClient
	if err = cont.Fill(rmqDef.DIClientPublisherEvents, &publisher); err != nil {
		return nil, err
	}

	var be receiver.Backend
	if be, err = receiver.NewBackend(
		emailRep,
		publisher,
		l,
	); err != nil {
		return nil, err
	}

	var smtpConfig receiver.ServerConfig
	if err = cfg.UnmarshalKey("smtp.server", &smtpConfig); err != nil {
		return nil, err
	}

	return receiver.NewSMTPReceiver(be, l, smtpConfig)
}
