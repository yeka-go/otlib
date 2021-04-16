package otsql

import "database/sql"

// Config tells if span should report table and/or query
type Config struct {
	WithTable  bool
	WithQuery  bool
	driverName string
}

func Open(driverName, dataSourceName string, cfg *Config) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return db, err
	}

	cg := Config{driverName: driverName}
	if cfg != nil {
		cg = *cfg
	}

	c := &connector{
		driver: mydriver{db.Driver()},
		dsn:    dataSourceName,
		cfg:    cg,
	}
	return sql.OpenDB(c), nil
}
