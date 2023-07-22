package rmq

// RMQ
//go:generate protoc -I=. -I=../../. -I=../../vendor --go-rmq_out=./../../pb --go-rmq_opt=paths=source_relative,gen_type=client publisher.proto
//go:generate protoc -I=. -I=../../. -I=../../vendor --go-rmq_out=./../../pb --go-rmq_opt=paths=source_relative,gen_type=server smtp_message_handler.proto

import (
	_ "github.com/requiemofthesouls/pigeomail/api/proto"

	_ "github.com/requiemofthesouls/pigeomail/pkg/tools/protoc-gen-go-rmq/proto"
)
