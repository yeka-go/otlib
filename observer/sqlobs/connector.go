package otsql

import (
	"context"
	"database/sql/driver"
)

type connector struct {
	dsn    string
	driver driver.Driver
	cfg    Config
}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.Driver().Open(c.dsn)
	if err != nil {
		return nil, err
	}

	return &connection{
		Conn: conn,
		ctx:  ctx,
		cfg:  c.cfg,
	}, nil
}

func (c *connector) Driver() driver.Driver {
	return c.driver
}
