package cmd

import (
	"fmt"

	"github.com/requiemofthesouls/pigeomail/pkg/modules/logger"
	logDef "github.com/requiemofthesouls/pigeomail/pkg/modules/logger/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/migrate"
	migratePostgresDef "github.com/requiemofthesouls/pigeomail/pkg/modules/migrate/postgres/def"
	"github.com/spf13/cobra"
)

var (
	migrationsPath string

	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate database",
	}

	createArgsValidator = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("migrate create need one arg - migration name (string)")
		}

		if args[0] == "" {
			return fmt.Errorf("migration name can't be empty")
		}

		return nil
	}

	migrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "apply all migrations",
		Long:  "Migrate the DB to the most recent version available",
		RunE:  migrateUpCmdHandler,
	}

	migrateDownCmd = &cobra.Command{
		Use:   "down",
		Short: "rollback all migrations",
		Long:  "Roll back all migrations",
		RunE:  migrateDownCmdHandler,
	}

	migrateDownByOneCmd = &cobra.Command{
		Use:   "down-by-one",
		Short: "rollback one transaction",
		Long:  "Roll back the version by 1",
		RunE:  migrateDownByOneCmdHandler,
	}

	migrateCreateCmd = &cobra.Command{
		Use:   "create [migration_name]",
		Short: "create migration",
		Long:  "Creates new migration file with the current timestamp",
		Args:  createArgsValidator,
		RunE:  migrateCreateCmdHandler,
	}
)

// Command init function.
func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateDownByOneCmd)
	migrateCmd.AddCommand(migrateCreateCmd)
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.PersistentFlags().StringVarP(&migrationsPath, "migrationsPath", "m", "migrations/pgsql", `path to migration files`)
}

func migrateUpCmdHandler(*cobra.Command, []string) (err error) {
	var m migrate.Migrate
	if err = diContainer.Fill(migratePostgresDef.DIWrapper, &m); err != nil {
		return err
	}

	m.SetPath(migrationsPath)

	var l logDef.Wrapper
	if err = diContainer.Fill(logDef.DIWrapper, &l); err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		return err
	}

	return nil
}

func migrateDownCmdHandler(*cobra.Command, []string) (err error) {
	var m migrate.Migrate
	if err = diContainer.Fill(migratePostgresDef.DIWrapper, &m); err != nil {
		return err
	}

	m.SetPath(migrationsPath)

	var l logger.Wrapper
	if err = diContainer.Fill(logDef.DIWrapper, &l); err != nil {
		return err
	}

	if err = m.Down(); err != nil {
		return err
	}

	return nil
}

func migrateDownByOneCmdHandler(*cobra.Command, []string) (err error) {
	var m migrate.Migrate
	if err = diContainer.Fill(migratePostgresDef.DIWrapper, &m); err != nil {
		return err
	}

	m.SetPath(migrationsPath)

	var l logger.Wrapper
	if err = diContainer.Fill(logDef.DIWrapper, &l); err != nil {
		return err
	}

	if err = m.DownByOne(); err != nil {
		return err
	}

	return nil
}

func migrateCreateCmdHandler(_ *cobra.Command, args []string) error {
	return migrate.CreateMigrationForPostgres(args[0], migrationsPath)
}
