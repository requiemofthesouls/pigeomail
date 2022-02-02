package telegram

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/mock/gomock"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"pigeomail/internal/adapters/rabbitmq"
	mocks2 "pigeomail/internal/adapters/rabbitmq/mock"
	"pigeomail/internal/domain/pigeomail"
	mocks "pigeomail/internal/domain/pigeomail/mock"
	"pigeomail/pkg/logger"
)

var (
	token  = "5123791453:AAHKpMEP5BB4L5jg691awNosKQGR8oD6oBc"
	domain = "shieldemail.ddns.net"
)

func TestBot_Run(t *testing.T) {
	type fields struct {
		api      *tgbotapi.BotAPI
		updates  tgbotapi.UpdatesChannel
		svc      pigeomail.Service
		consumer rabbitmq.Consumer
		domain   string
		logger   *logr.Logger
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				api:      tt.fields.api,
				updates:  tt.fields.updates,
				svc:      tt.fields.svc,
				consumer: tt.fields.consumer,
				domain:   tt.fields.domain,
				logger:   tt.fields.logger,
			}
			b.Run()
		})
	}
}

func TestBot_handleCommand(t *testing.T) {
	type fields struct {
		api      *tgbotapi.BotAPI
		updates  tgbotapi.UpdatesChannel
		svc      pigeomail.Service
		consumer rabbitmq.Consumer
		domain   string
		logger   *logr.Logger
	}
	type args struct {
		update *tgbotapi.Update
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				api:      tt.fields.api,
				updates:  tt.fields.updates,
				svc:      tt.fields.svc,
				consumer: tt.fields.consumer,
				domain:   tt.fields.domain,
				logger:   tt.fields.logger,
			}
			b.handleCommand(tt.args.update)
		})
	}
}

func TestBot_handleError(t *testing.T) {
	type fields struct {
		api      *tgbotapi.BotAPI
		updates  tgbotapi.UpdatesChannel
		svc      pigeomail.Service
		consumer rabbitmq.Consumer
		domain   string
		logger   *logr.Logger
	}
	type args struct {
		err    error
		chatID int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				api:      tt.fields.api,
				updates:  tt.fields.updates,
				svc:      tt.fields.svc,
				consumer: tt.fields.consumer,
				domain:   tt.fields.domain,
				logger:   tt.fields.logger,
			}
			b.handleError(tt.args.err, tt.args.chatID)
		})
	}
}

func TestBot_incomingEmailConsumer(t *testing.T) {
	type fields struct {
		api      *tgbotapi.BotAPI
		updates  tgbotapi.UpdatesChannel
		svc      pigeomail.Service
		consumer rabbitmq.Consumer
		domain   string
		logger   *logr.Logger
	}
	type args struct {
		msg *amqp.Delivery
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				api:      tt.fields.api,
				updates:  tt.fields.updates,
				svc:      tt.fields.svc,
				consumer: tt.fields.consumer,
				domain:   tt.fields.domain,
				logger:   tt.fields.logger,
			}
			b.incomingEmailConsumer(tt.args.msg)
		})
	}
}

func TestBot_runBot(t *testing.T) {
	type fields struct {
		api      *tgbotapi.BotAPI
		updates  tgbotapi.UpdatesChannel
		svc      pigeomail.Service
		consumer rabbitmq.Consumer
		domain   string
		logger   *logr.Logger
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				api:      tt.fields.api,
				updates:  tt.fields.updates,
				svc:      tt.fields.svc,
				consumer: tt.fields.consumer,
				domain:   tt.fields.domain,
				logger:   tt.fields.logger,
			}
			b.runBot()
		})
	}
}

func TestBot_runConsumer(t *testing.T) {
	type fields struct {
		api      *tgbotapi.BotAPI
		updates  tgbotapi.UpdatesChannel
		svc      pigeomail.Service
		consumer rabbitmq.Consumer
		domain   string
		logger   *logr.Logger
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bot{
				api:      tt.fields.api,
				updates:  tt.fields.updates,
				svc:      tt.fields.svc,
				consumer: tt.fields.consumer,
				domain:   tt.fields.domain,
				logger:   tt.fields.logger,
			}
			b.runConsumer()
		})
	}
}

func TestNewBot(t *testing.T) {
	logger.Init()

	var (
		ctrl = gomock.NewController(t)
		ctx  = context.Background()

		svc  = mocks.NewMockService(ctrl)
		cons = mocks2.NewMockConsumer(ctrl)
		l    = logger.GetLogger()
	)

	type args struct {
		ctx         context.Context
		debug       bool
		webhookMode bool
		token       string
		domain      string
		port        string
		cert        string
		key         string
		svc         pigeomail.Service
		cons        rabbitmq.Consumer
	}
	tests := []struct {
		name        string
		args        args
		wantBot     *Bot
		wantErr     bool
		expectedErr string
	}{
		{
			name: "new updates bot",
			args: args{
				ctx:         ctx,
				debug:       true,
				webhookMode: false,
				token:       token,
				domain:      domain,
				svc:         svc,
				cons:        cons,
			},
			wantBot: &Bot{
				svc:      svc,
				consumer: cons,
				domain:   domain,
				logger:   l,
			},
			wantErr: false,
		},
		{
			name: "bad token",
			args: args{
				ctx:         ctx,
				debug:       true,
				webhookMode: false,
				token:       "badToken",
				domain:      domain,
				svc:         svc,
				cons:        cons,
			},
			wantBot: &Bot{
				svc:      svc,
				consumer: cons,
				domain:   domain,
				logger:   l,
			},
			wantErr:     true,
			expectedErr: "Not Found",
		},
		{
			name: "bad cert",
			args: args{
				ctx:         ctx,
				debug:       true,
				webhookMode: true,
				token:       token,
				domain:      domain,
				port:        "8443",
				cert:        "badCertFile",
				key:         "",
				svc:         svc,
				cons:        cons,
			},
			wantBot: &Bot{
				svc:      svc,
				consumer: cons,
				domain:   domain,
				logger:   l,
			},
			wantErr:     true,
			expectedErr: "Post \"https://api.telegram.org/bot5123791453:AAHKpMEP5BB4L5jg691awNosKQGR8oD6oBc/setWebhook\": open badCertFile: no such file or directory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBot, gotErr := NewBot(tt.args.ctx, tt.args.debug, tt.args.webhookMode, tt.args.token, tt.args.domain, tt.args.port, tt.args.cert, tt.args.key, tt.args.svc, tt.args.cons)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("NewBot() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			if tt.wantErr {
				assert.EqualError(t, gotErr, tt.expectedErr)
				return
			}

			assert.Equal(t, gotBot.svc, tt.wantBot.svc)
			assert.Equal(t, gotBot.consumer, tt.wantBot.consumer)
			assert.Equal(t, gotBot.domain, tt.wantBot.domain)
			assert.Equal(t, gotBot.logger, tt.wantBot.logger)

		})
	}
}

func Test_getUpdatesChan(t *testing.T) {
	type args struct {
		log   *logr.Logger
		tgAPI *tgbotapi.BotAPI
	}
	tests := []struct {
		name        string
		args        args
		wantUpdates tgbotapi.UpdatesChannel
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUpdates := getUpdatesChan(tt.args.log, tt.args.tgAPI); !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
				t.Errorf("getUpdatesChan() = %v, want %v", gotUpdates, tt.wantUpdates)
			}
		})
	}
}

func Test_getWebhookUpdatesChan(t *testing.T) {
	type args struct {
		tgAPI  *tgbotapi.BotAPI
		domain string
		port   string
		cert   string
		key    string
	}
	tests := []struct {
		name        string
		args        args
		wantUpdates tgbotapi.UpdatesChannel
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUpdates, err := getWebhookUpdatesChan(tt.args.tgAPI, tt.args.domain, tt.args.port, tt.args.cert, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getWebhookUpdatesChan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUpdates, tt.wantUpdates) {
				t.Errorf("getWebhookUpdatesChan() gotUpdates = %v, want %v", gotUpdates, tt.wantUpdates)
			}
		})
	}
}
