package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Handler - func, which should be executed within transaction
type Handler func(ctx context.Context) error

type Client interface {
	DB() DB
	Close() error
}

// TxManager executes handler, which defined by user within certain tx
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Query is a request wrapper, which request
type Query struct {
	Name string
	QueryRaw string
}

// Transactor is an interface for working with txs
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// SQLExecer combines NamedExeceer and QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer is interface for working with named requests with tags in structs
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecer is an interface for working with ordinal requests
type QueryExecer interface {
	ExecContect(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger is an interface for working with DB connection
type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}