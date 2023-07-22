module github.com/requiemofthesouls/pigeomail

go 1.20

replace (
	github.com/requiemofthesouls/pigeomail/api => ./api
	//github.com/requiemofthesouls/pigeomail/pkg/tools/protoc-gen-go-rmq => ./pkg/tools/protoc-gen-go-rmq
)

require (
	github.com/emersion/go-smtp v0.16.0
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/jackc/pgx/v5 v5.3.1
	github.com/jhillyerd/enmime v0.11.0
	github.com/looplab/fsm v1.0.1
	github.com/pressly/goose v2.7.0+incompatible
	github.com/rabbitmq/amqp091-go v1.8.1
	github.com/requiemofthesouls/pigeomail/api v0.0.0-00010101000000-000000000000
	github.com/requiemofthesouls/pigeomail/pkg/tools/protoc-gen-go-rmq v0.0.0-00010101000000-000000000000
	github.com/sarulabs/di/v2 v2.4.2
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.15.0
	github.com/stretchr/testify v1.8.1
	go.uber.org/zap v1.24.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/cention-sany/utf7 v0.0.0-20170124080048-26cad61bd60a // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emersion/go-sasl v0.0.0-20220912192320-0145f2c60ead // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gogs/chardet v0.0.0-20211120154057-b7413eaefb8f // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.0 // indirect
	github.com/jaytaylor/html2text v0.0.0-20211105163654-bc68cce691ba // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.7 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.10.0 // indirect
	golang.org/x/net v0.11.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.10.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230626202813-9b080da550b3 // indirect
	google.golang.org/grpc v1.56.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
