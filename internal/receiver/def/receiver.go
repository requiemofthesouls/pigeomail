package def

import (
	"github.com/requiemofthesouls/pigeomail/internal/receiver"
	repDef "github.com/requiemofthesouls/pigeomail/internal/repository/def"
	cfgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/config/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
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

	//var publisher rmqDef.Publisher
	//if err = cont.Fill(rmqDef.DIAMQPPublisher, &publisher); err != nil {
	//	return nil, err
	//}

	var be receiver.Backend
	if be, err = receiver.NewBackend(
		emailRep,
		//publisher,
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
