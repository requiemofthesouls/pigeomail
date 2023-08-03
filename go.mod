module github.com/requiemofthesouls/pigeomail

go 1.20

replace github.com/requiemofthesouls/pigeomail/api => ./api

require (
	github.com/emersion/go-smtp v0.17.0
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/jhillyerd/enmime v1.0.0
	github.com/looplab/fsm v1.0.1
	github.com/r3labs/sse/v2 v2.10.0
	github.com/requiemofthesouls/config v0.0.1
	github.com/requiemofthesouls/container v0.0.1
	github.com/requiemofthesouls/logger v0.0.2
	github.com/requiemofthesouls/migrate v0.0.2
	github.com/requiemofthesouls/monitoring v0.0.1
	github.com/requiemofthesouls/pigeomail/api v0.0.0-00010101000000-000000000000
	github.com/requiemofthesouls/postgres v0.0.2
	github.com/requiemofthesouls/protoc-gen-go-rmq v0.0.0-20230729142903-06f69ddd1b7c
	github.com/requiemofthesouls/svc-grpc v0.0.5
	github.com/requiemofthesouls/svc-http v0.0.1
	github.com/requiemofthesouls/svc-rmq v0.0.1
	github.com/spf13/cobra v1.7.0
	go.uber.org/zap v1.25.0
	google.golang.org/grpc v1.57.0
)

require (
	github.com/VictoriaMetrics/metrics v1.24.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cention-sany/utf7 v0.0.0-20170124080048-26cad61bd60a // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emersion/go-sasl v0.0.0-20220912192320-0145f2c60ead // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/gogs/chardet v0.0.0-20211120154057-b7413eaefb8f // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.2 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jaytaylor/html2text v0.0.0-20230321000545-74c2419ad056 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/pressly/goose v2.7.0+incompatible // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/rabbitmq/amqp091-go v1.8.1 // indirect
	github.com/requiemofthesouls/client-errors v0.0.1 // indirect
	github.com/requiemofthesouls/user-client v0.0.1 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/sarulabs/di/v2 v2.4.2 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.16.0 // indirect
	github.com/ssor/bom v0.0.0-20170718123548-6386211fdfcf // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/valyala/fastrand v1.1.0 // indirect
	github.com/valyala/histogram v1.2.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/net v0.13.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	google.golang.org/genproto v0.0.0-20230731193218-e0aa005b6bdf // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230731193218-e0aa005b6bdf // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230731193218-e0aa005b6bdf // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
