package def

import (
	"context"

	cfgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/config/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/postgres"
)

const (
	DIWrapper      = "postgres.wrapper"
	DIWrapperSqlDB = "postgres.wrapper.sql_db"
)

type (
	Wrapper      = postgres.Wrapper
	WrapperSqlDB = *postgres.SqlDB
)

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(
			container.Def{
				Name: DIWrapper,
				Build: func(container container.Container) (interface{}, error) {
					var cfg cfgDef.Wrapper
					if err := container.Fill(cfgDef.DIWrapper, &cfg); err != nil {
						return nil, err
					}

					var pgCfg postgres.Config
					if err := cfg.UnmarshalKey("postgres", &pgCfg); err != nil {
						return nil, err
					}

					return postgres.New(context.Background(), pgCfg)
				},
				Close: func(obj interface{}) error {
					obj.(postgres.Wrapper).Close()
					return nil
				},
			},
			container.Def{
				Name: DIWrapperSqlDB,
				Build: func(container container.Container) (interface{}, error) {
					var pgWrapper postgres.Wrapper
					if err := container.Fill(DIWrapper, &pgWrapper); err != nil {
						return nil, err
					}

					return postgres.NewSqlDB(pgWrapper)
				},
				Close: func(obj interface{}) error {
					_ = obj.(*postgres.SqlDB).Close()
					return nil
				},
			})
	})
}
