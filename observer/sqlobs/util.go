package otsql

import (
	"github.com/opentracing/opentracing-go/log"
)

func AttrFromQuery(cfg Config, query string) []log.Field {
	att := make([]log.Field, 0)

	op, table, err := Parse(query)
	if err == nil {
		att = append(att, log.String("db.operator", op))
	} else {
		att = append(att, log.String("db.operator", err.Error()))
	}

	if cfg.WithTable {
		att = append(att, log.String("db.table", table))
	}

	if cfg.WithQuery {
		att = append(att, log.String("db.query", query))
	}
	return att
}

func AttrFromConfig(cfg Config) []log.Field {
	att := []log.Field{
		log.String("db.drivername", cfg.driverName),
	}
	return att
}
