package def

import (
	"github.com/requiemofthesouls/container"
	"github.com/requiemofthesouls/pigeomail/internal/repository"
	pgDef "github.com/requiemofthesouls/postgres/def"
)

const (
	DIDBRepositoryTelegramUsers = "db.repository.telegram_users"
)

type (
	TelegramUsers = repository.TelegramUsers
)

var dbDeps = map[string]func(pgDef.Wrapper) interface{}{
	DIDBRepositoryTelegramUsers: func(db pgDef.Wrapper) interface{} { return repository.NewUsers(db) },
}

func init() {
	var defs = make([]container.Def, 0, len(dbDeps))
	for defDB, fn := range dbDeps {
		var fnRepo = fn

		defs = append(defs, container.Def{
			Name: defDB,
			Build: func(cont container.Container) (interface{}, error) {
				var db pgDef.Wrapper
				if err := cont.Fill(pgDef.DIWrapper, &db); err != nil {
					return nil, err
				}

				return fnRepo(db), nil
			},
		})
	}

	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(defs...)
	})
}
