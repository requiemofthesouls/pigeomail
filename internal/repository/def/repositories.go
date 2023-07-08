package def

import (
	"github.com/requiemofthesouls/pigeomail/internal/repository"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	pgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/postgres/def"
	"github.com/requiemofthesouls/pigeomail/pkg/state"
)

const (
	DIDBRepositoryTelegramUsers          = "db.repository.telegram_users"
	DIDBRepositoryTelegramUsersWithState = "db.repository.telegram_users_with_state"
)

type (
	TelegramUsers          = repository.TelegramUsers
	TelegramUsersWithState = repository.TelegramUsersWithState
)

var dbDeps = map[string]func(pgDef.Wrapper) interface{}{
	DIDBRepositoryTelegramUsers: func(db pgDef.Wrapper) interface{} { return repository.NewUsers(db) },
	DIDBRepositoryTelegramUsersWithState: func(db pgDef.Wrapper) interface{} {
		return repository.NewUsersWithState(repository.NewUsers(db), state.NewState())
	},
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
