package cmd

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/requiemofthesouls/pigeomail/internal/telegram/def"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "tg_bot",
		Short: "start telegram bot",
		RunE:  startTGBot,
	})
}

func startTGBot(_ *cobra.Command, _ []string) error {
	var l logDef.Wrapper
	if err := diContainer.Fill(logDef.DIWrapper, &l); err != nil {
		return err
	}

	var tgBot *def.TGBot
	if err := diContainer.Fill(def.DITGBot, &tgBot); err != nil {
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
	go func() { defer wg.Done(); tgBot.Run(ctx, l) }()

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
