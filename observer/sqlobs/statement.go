package otsql

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// - driver.Stmt
// - driver.StmtExecContext
// - driver.StmtQueryContext
type statement struct {
	driver.Stmt
	ctx  context.Context
	span opentracing.Span
	cfg  Config
}

// Close implements driver.Stmt
func (st *statement) Close() error {
	defer st.span.Finish()
	return st.Stmt.Close()
}

// Exec implements driver.Stmt
func (st *statement) Exec(args []driver.Value) (driver.Result, error) {
	span, _ := opentracing.StartSpanFromContext(st.ctx, "db.statement.exec")
	defer span.Finish()

	return st.Stmt.Exec(args)
}

// Query implements driver.Stmt
func (st *statement) Query(args []driver.Value) (driver.Rows, error) {
	span, _ := opentracing.StartSpanFromContext(st.ctx, "db.statement.query")

	r, err := st.Stmt.Query(args)
	if err != nil {
		span.LogFields(log.String("error", err.Error()))
		span.Finish()
	}
	return &rows{r, span}, err
}

// ExecContext implements driver.StmtExecContext
func (st *statement) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(st.ctx, "db.statement.exec_context")
	defer span.Finish()

	ec, ok := st.Stmt.(driver.StmtExecContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	return ec.ExecContext(ctx, args)
}

// QueryContext implements driver.StmtQueryContext
func (st *statement) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(st.ctx, "db.statement.query_context")

	qc, ok := st.Stmt.(driver.StmtQueryContext)
	if !ok {
		span.LogFields(log.String("error", driver.ErrSkip.Error()))
		span.Finish()
		return nil, driver.ErrSkip
	}

	r, err := qc.QueryContext(ctx, args)
	if err != nil {
		span.LogFields(log.String("error", err.Error()))
		span.Finish()
	}
	return &rows{r, span}, err
}
