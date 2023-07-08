package postgres

//go:generate mockery --name=Tx --inpackage --output=. --filename=tx_mock.go --structname=TxMock

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	CommandTag = pgconn.CommandTag

	Rows = pgx.Rows
	Row  = pgx.Row

	Tx = pgx.Tx
)
