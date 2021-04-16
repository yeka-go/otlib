package otsql

import (
	"database/sql/driver"

	"github.com/opentracing/opentracing-go"
)

type rows struct {
	driver.Rows
	span opentracing.Span
}

func (r *rows) Close() error {
	defer r.span.Finish()
	return r.Rows.Close()
}
