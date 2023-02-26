package cmd

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/spf13/cobra"

	"github.com/requiemofthesouls/pigeomail/internal/receiver/def"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "receiver",
		Short: "start receiver service",
		RunE:  startReceiver,
	})
}

func startReceiver(_ *cobra.Command, _ []string) error {
	var receiver *def.Receiver
	if err := diContainer.Fill(def.DISMTPReceiver, &receiver); err != nil {
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
