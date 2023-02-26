package def

import (
	"github.com/requiemofthesouls/pigeomail/internal/repository"
	cfgDef "github.com/requiemofthesouls/pigeomail/pkg/modules/config/def"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/container"
	"github.com/requiemofthesouls/pigeomail/pkg/modules/mongo"
	mongoDef "github.com/requiemofthesouls/pigeomail/pkg/modules/mongo/def"
	"github.com/requiemofthesouls/pigeomail/pkg/state"
)

const DIDBRepositoryEmail = "db.repository.email"
const DIDBRepositoryEmailState = "db.repository.email_state"

type Email = repository.Email
type EmailState = repository.EmailState

func initEmailRepo() container.Def {
	return container.Def{
		Name: DIDBRepositoryEmail,
		Build: func(cont container.Container) (_ interface{}, err error) {
			var cfg cfgDef.Wrapper
			if err = cont.Fill(cfgDef.DIWrapper, &cfg); err != nil {
				return nil, err
			}

			var mCfg mongo.Config
			if err = cfg.UnmarshalKey("mongo", &mCfg); err != nil {
				return nil, err
			}

			var db mongo.Wrapper
			if err = cont.Fill(mongoDef.DIWrapper, &db); err != nil {
				return nil, err
			}

			var repo = repository.NewEmail(db)

			return repo, nil
		},
	}
}

func initEmailStateRepo() container.Def {
	return container.Def{
		Name: DIDBRepositoryEmailState,
		Build: func(cont container.Container) (_ interface{}, err error) {
			var repo Email
			if err = cont.Fill(DIDBRepositoryEmail, &repo); err != nil {
				return nil, err
			}

			var st = state.NewState()
			var stateRepo = repository.NewEmailState(repo, st)

			return stateRepo, nil
		},
	}
}

func init() {
	container.Register(func(builder *container.Builder, _ map[string]interface{}) error {
		return builder.Add(initEmailRepo(), initEmailStateRepo())
	})
}
