package otsql

import (
	"database/sql/driver"

	"github.com/opentracing/opentracing-go"
)

type tx struct {
	driver.Tx
	ctxClearer func()
	span       opentracing.Span
}

func (tx *tx) Commit() error {
	tx.span.SetOperationName("db.transaction commit")
	defer tx.span.Finish()
	defer tx.ctxClearer()
	return tx.Tx.Commit()
}

func (tx *tx) Rollback() error {
	tx.span.SetOperationName("db.transaction rollback")
	defer tx.span.Finish()
	defer tx.ctxClearer()
	return tx.Tx.Rollback()
}
