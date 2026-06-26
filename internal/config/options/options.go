package options

import "time"

type PostgresOptions struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
}

type HTTPOptions struct {
	Addr  string
	Bind  string
	Token string
}
