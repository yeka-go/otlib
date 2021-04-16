package otsql

import (
	"context"
	"database/sql/driver"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// - driver.Conn
// x driver.ConnBeginTx
// x driver.ConnPrepareContext
// - driver.Pinger
// x driver.Execer
// x driver.ExecerContext
// - driver.Queryer
// - driver.QueryerContext
type connection struct {
	driver.Conn
	ctx   context.Context
	txctx context.Context
	cfg   Config
}

// Close implements driver.Conn
func (c *connection) Close() error {
	return c.Conn.Close()
}

// Prepare implements driver.Conn
func (c *connection) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(c.ctx, query)
}

// Begin implements driver.Conn
func (c *connection) Begin() (driver.Tx, error) {
	return c.BeginTx(c.ctx, driver.TxOptions{})
}

// PrepareContext implements driver.ConnPrepareContext
func (c *connection) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if c.txctx != nil {
		ctx = c.txctx
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, "db.prepare_context")
	span.LogFields(AttrFromConfig(c.cfg)...)
	span.LogFields(AttrFromQuery(c.cfg, query)...)

	cn, ok := c.Conn.(driver.ConnPrepareContext)
	var st driver.Stmt
	var err error
	if ok {
		st, err = cn.PrepareContext(ctx, query)
	} else {
		st, err = c.Conn.Prepare(query)
	}

	if err != nil {
		span.LogFields(log.String("error", err.Error()))
		span.Finish()
		return st, err
	}

	return &statement{st, ctx, span, c.cfg}, err
}

// BeginTx implements driver.ConnBeginTx
func (c *connection) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "db.transaction")
	span.LogFields(AttrFromConfig(c.cfg)...)

	c.txctx = ctx
	ctxClearer := func() {
		c.txctx = nil
	}

	var t driver.Tx
	var err error
	bt, ok := c.Conn.(driver.ConnBeginTx)
	if ok {
		t, err = bt.BeginTx(ctx, opts)
	} else {
		t, err = c.Conn.Begin()
	}

	if err != nil {
		span.LogFields(log.String("error", err.Error()))
		span.Finish()
		return t, err
	}

	return &tx{t, ctxClearer, span}, nil
}

// Ping implements driver.Pinger
func (c *connection) Ping(ctx context.Context) error {
	p, ok := c.Conn.(driver.Pinger)
	if !ok {
		return driver.ErrSkip
	}
	return p.Ping(ctx)
}

// QueryContext implements driver.QueryerContext
func (c *connection) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if c.txctx != nil {
		ctx = c.txctx
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, "db.query_context")
	span.LogFields(AttrFromConfig(c.cfg)...)

	qc, ok := c.Conn.(driver.QueryerContext)
	if !ok {
		span.LogFields(log.String("error", driver.ErrSkip.Error()))
		span.Finish()
		return nil, driver.ErrSkip
	}

	span.LogFields(AttrFromQuery(c.cfg, query)...)

	r, err := qc.QueryContext(ctx, query, args)
	if err != nil {
		span.LogFields(log.String("error", err.Error()))
		span.Finish()
		return r, err
	}
	return &rows{r, span}, nil
}

// ExecContext implements driver.ExecerContext
func (c *connection) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if c.txctx != nil {
		ctx = c.txctx
	}
	span, ctx := opentracing.StartSpanFromContext(ctx, "db.exec_context")
	defer span.Finish()
	span.LogFields(AttrFromConfig(c.cfg)...)

	qe, ok := c.Conn.(driver.ExecerContext)
	if !ok {
		span.LogFields(log.String("error", driver.ErrSkip.Error()))
		return nil, driver.ErrSkip
	}

	span.LogFields(AttrFromQuery(c.cfg, query)...)

	res, err := qe.ExecContext(ctx, query, args)
	if err != nil {
		span.LogFields(log.String("error", err.Error()))
	}
	return res, err
}
