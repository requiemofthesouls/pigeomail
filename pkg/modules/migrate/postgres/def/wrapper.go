package def

import (
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/migrate"
	migratePg "github.com/requiemofthesouls/pigeomail/pkg/modules/migrate/postgres"
	pgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/postgres/def"
)

const DIWrapper = "migrate.postgres"

type Wrapper = migrate.Migrate

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DIWrapper,
			Build: func(cont container.Container) (interface{}, error) {
				var db pgDef.WrapperSqlDB
				if err := cont.Fill(pgDef.DIWrapperSqlDB, &db); err != nil {
					return nil, err
				}

				return migratePg.NewWrapper(db), nil
			},
		})
	})
}
