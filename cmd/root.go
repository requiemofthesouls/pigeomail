package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/requiemofthesouls/container"
)

var (
	diContainer      container.Container
	cfgPath          string
	stopNotification = make(chan struct{})

	// version
	commitHash = "0000000000000000000000000000000000000000"
	branch     string
	tag        = "v0.0.0"
	buildDate  string
	builtBy    string

	// Root command.
	rootCmd = &cobra.Command{
		Use:           "pigeomail [command]",
		Long:          "pigeomail project",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			if diContainer, err = container.Instance(
				[]string{
					container.App,
					container.Request,
					container.SubRequest,
				},
				map[string]interface{}{
					"config":      cfgPath,
					"cli_cmd":     cmd,
					"cli_args":    args,
					"commit_hash": commitHash,
					"branch":      branch,
					"tag":         tag,
					"build_date":  buildDate,
					"built_by":    builtBy,
				}); err != nil {
				return err
			}

			// graceful stop
			go func() {
				var c = make(chan os.Signal, 1)
				signal.Notify(c,
					syscall.SIGHUP,
					syscall.SIGINT,
					syscall.SIGTERM,
				)

				<-c

				stopNotification <- struct{}{}
			}()

			return err
		},
	}
)

func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&cfgPath, "config", "c", "config.yaml", "config file")
	if err := rootCmd.Execute(); err != nil {
		return err
	}

	if diContainer != nil {
		return diContainer.Delete()
	}

	return nil
}
