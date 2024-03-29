package cmd

import (
	"context"
	"errors"
	"sync"
	"time"

	logDef "github.com/requiemofthesouls/logger/def"
	receiverDef "github.com/requiemofthesouls/pigeomail/internal/receiver/def"
	tgBotDef "github.com/requiemofthesouls/pigeomail/internal/telegram/def"
	grpcService "github.com/requiemofthesouls/svc-grpc"
	grpcServiceDef "github.com/requiemofthesouls/svc-grpc/def"
	httpService "github.com/requiemofthesouls/svc-http"
	httpServiceDef "github.com/requiemofthesouls/svc-http/def"
	rmqDef "github.com/requiemofthesouls/svc-rmq/def"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "start pigeomail service",
		RunE:  startService,
	})
}

func startService(_ *cobra.Command, _ []string) error {
	var l logDef.Wrapper
	if err := diContainer.Fill(logDef.DIWrapper, &l); err != nil {
		return err
	}

	var grpcServerManager grpcService.Manager
	if err := diContainer.Fill(grpcServiceDef.DIServerManager, &grpcServerManager); err != nil {
		return err
	}

	var httpServerManager httpService.Manager
	if err := diContainer.Fill(httpServiceDef.DIServerManager, &httpServerManager); err != nil {
		return err
	}

	var tgBot *tgBotDef.Bot
	if err := diContainer.Fill(tgBotDef.DITelegramBot, &tgBot); err != nil {
		return err
	}

	var rmqManager rmqDef.Manager
	if err := diContainer.Fill(rmqDef.DIManager, &rmqManager); err != nil {
		return err
	}

	var receiver *receiverDef.Receiver
	if err := diContainer.Fill(receiverDef.DISMTPReceiver, &receiver); err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	// graceful shutdown
	go func() {
		defer cancel()
		<-stopNotification
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); grpcServerManager.StartAll(ctx) }()

	wg.Add(1)
	go func() { defer wg.Done(); httpServerManager.StartAll(ctx) }()

	wg.Add(1)
	go func() { defer wg.Done(); tgBot.Start(ctx) }()

	wg.Add(1)
	go func() { defer wg.Done(); rmqManager.StartAll(ctx) }()

	wg.Add(1)
	go func() { defer wg.Done(); receiver.Run(ctx) }()

	<-ctx.Done()

	var waitChan = make(chan struct{})
	go func() {
		wg.Wait()
		waitChan <- struct{}{}
	}()

	select {
	case <-time.After(time.Second * 5):
		return errors.New("couldn't stop service within the specified timeout (5 sec)")
	case <-waitChan:
		return nil
	}
}
