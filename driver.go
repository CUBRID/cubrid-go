// Copyright 2019 The Go-CUBRID-Driver Authors. All rights reserved.
//
// Package cubrid provides a CUBRID driver for Go's database/sql package.
//
// The driver should be used via the database/sql package:
//
//  import "database/sql"
//  import _ "github.com/CUBRID/cubrid-go"

package cubrid

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

// CubridDriver is exported to make the driver directly accessible.
// In general the driver is used via the database/sql package.
type CubridDriver struct{}

// Open new Connection.
func (d CubridDriver) Open(conn_url string) (driver.Conn, error) {

	c := &connector {
		url: &conn_url,
	}

	return c.Connect(context.Background())
}

func (d CubridDriver) OpenConnector(conn_url string) (driver.Connector, error) {
	var err error

	err = nil

	c := &connector {
		url: &conn_url,
	}

	return c, err
}

func init() {
	sql.Register("cubrid", &CubridDriver{})
}
