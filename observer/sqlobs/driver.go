package otsql

import (
	"database/sql/driver"
)

type mydriver struct {
	driver driver.Driver
}

func (md mydriver) Open(name string) (driver.Conn, error) {
	//log.Println("opening connection:", name)
	return md.driver.Open(name)
}
