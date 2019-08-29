// Copyright 2019 The Go-CUBRID-Driver Authors. All rights reserved.
//
// Package cubrid provides a CUBRID driver for Go's database/sql package.
//
// The driver should be used via the database/sql package:
//
//  import "database/sql"
//  import _ "github.com/CUBRID/cubrid-go"

package cubrid

// #include <cas_cci.h>
import "C"

import (
	"context"
	"database/sql/driver"
	"errors"
)

type connector struct {
	url	*string
}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	var err		error
	var err_buf	C.T_CCI_ERROR

	conn_url := C.CString(*c.url)

	res := C.cci_connect_with_url_ex(conn_url, nil, nil, &err_buf)

	h := &cub_conn {
		handle:	int(res),
		cancel:	false,
		closed:	true,
	}

	if (res < 0) {
		err = errors.New(C.GoString(&err_buf.err_msg[0]))
		return h, err
	}

	return h, nil
}

func (c *connector) Driver() driver.Driver {
	return &CubridDriver{}
}

