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
	"database/sql/driver"
	"errors"
	"time"
)

type cub_conn struct {
	handle		int
	tr_start	time.Time
	closed		bool
	cancel		bool
}

func (c *cub_conn) Commit() (error) {
	var err error
	var err_buf C.T_CCI_ERROR

	if c.handle < 0 {
		return driver.ErrBadConn
	}

	res := C.cci_end_tran(C.int(c.handle), 1, &err_buf)
	if res < 0 {
		err = errors.New(C.GoString(&err_buf.err_msg[0]))
		return err
	}

	return nil
}

func (c *cub_conn) Rollback() (error) {
	var err error
	var err_buf C.T_CCI_ERROR

	if c.handle < 0 {
		return driver.ErrBadConn
	}

	res := C.cci_end_tran(C.int(c.handle), 2, &err_buf)
	if res < 0 {
		err = errors.New(C.GoString(&err_buf.err_msg[0]))
	}

	return err
}

func (c *cub_conn) Prepare(query string) (driver.Stmt, error) {
	var err error
	var flag C.char
	var err_buf C.T_CCI_ERROR
	var res C.int

	err = nil

	if c.handle < 0 {
		return nil, driver.ErrBadConn
	}

	sql := C.CString(query)
	res = C.cci_prepare(C.int(c.handle), sql, flag, &err_buf)

	if res < 0 {
		err := errors.New(C.GoString(&err_buf.err_msg[0]))
		return nil, err
	}

	stmt := &cub_stmt {
		conn: c,
		handle: -1,
		params: 0,
	}

	stmt.handle = int(res)

	res = C.cci_get_bind_num(res)
	if res < 0 {
		err := errors.New("connect: cci_get_bind_num, bind param numbers error")
		return stmt, err
	}

	stmt.params = int(res)

	return stmt, err
}

func (c *cub_conn) Begin() (driver.Tx, error) {
	if c.handle < 0 {
		return c, driver.ErrBadConn
	}

	c.tr_start = time.Now()
	return c, nil
}

func (c* cub_conn) Close() (error) {
	var err_buf C.T_CCI_ERROR

	if c.handle < 0 {
		return driver.ErrBadConn
	}

	res := C.cci_disconnect(C.int(c.handle), &err_buf)

	if res < 0 {
		err := errors.New(C.GoString(&err_buf.err_msg[0]))
		return err
	}

	c.closed = true

	return nil
}
