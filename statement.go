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
// #include <col_info.h>
import "C"

import (
	"database/sql/driver"
	"errors"
	"reflect"
	"unsafe"
)

type cub_stmt struct {
	conn	*cub_conn
	handle	int
	params	int
	flag	byte
}

func (s *cub_stmt) Close() error {
	if s.conn == nil || s.conn.closed || s.conn.handle < 0 {
		return driver.ErrBadConn
	}

	res := C.cci_close_req_handle(C.int(s.handle))
	if res < 0 {
		return errors.New("Close: bad statement handle")
	}

	s.handle = -1
	s.params = 0

	return  nil
}

func (s *cub_stmt) NumInput() int {
	return s.params
}

func (s *cub_stmt) Exec(args []driver.Value) (driver.Result, error) {
	var handle	C.int
	var flag	C.char
	var err_buf	C.T_CCI_ERROR

	handle = C.int(s.handle)
	flag = C.char(s.flag)

	res := C.cci_execute(handle, flag, 0, &err_buf)
	if int(res) < 0 {
		err := errors.New(C.GoString(&err_buf.err_msg[0]))
		return nil, err
	}

	return &cub_result {
		affected_rows: int64(res),
		last_insert: 0,
	}, nil
}

func (s *cub_stmt) Query(args []driver.Value) (driver.Rows, error) {
	var handle	C.int
	var flag	C.char
	var err_buf	C.T_CCI_ERROR
	var col_nums	C.int
	var col_info*	C.T_CCI_COL_INFO
	var ci		[]C.T_CCI_COLUMN_INFO
	var stmt_type	C.T_CCI_CUBRID_STMT

	handle = C.int(s.handle)
	flag = C.char(s.flag)

	col_info = C.cci_get_result_info(handle, &stmt_type, &col_nums)

	hdr := reflect.SliceHeader {
		Data:	uintptr(unsafe.Pointer(col_info)),
		Len:	C.sizeof_T_CCI_COL_INFO,
		Cap:	C.sizeof_T_CCI_COL_INFO,
	}

	ci = *(*[]C.T_CCI_COLUMN_INFO)(unsafe.Pointer(&hdr))

	if ci != nil {
		col_list := make([]cub_col, col_nums)
		name_list := make([]string, col_nums)
		for i:= 0; i < int(col_nums); i++ {
			name_list[i] = C.GoString(ci[i].col_name)
			col_list[i].table = C.GoString(ci[i].class_name)
			col_list[i].col_type = byte(ci[i].ext_type)
			col_list[i].scale = int(ci[i].scale)
			col_list[i].precision = int(ci[i].precision)
		}

		rs := &cub_result {
			columns:	col_list,
			names:		name_list,
			affected_rows:	0,
			last_insert:	0,
		}

		rows := &cub_rows {
			conn: s.conn.handle,
			handle: s.handle,
			result: rs,
		}

		res := C.cci_execute(handle, flag, 0, &err_buf)
		if int(res) < 0 {
			err := errors.New(C.GoString(&err_buf.err_msg[0]))
			return nil, err
		}
		rows.result.affected_rows = int64(res)

		res = C.cci_fetch_size(handle, 100)
		if int(res) < 0 {
			err := errors.New("Query: cci_fetch_size error")
			return nil, err
		}

		return rows, nil
	}

	return nil, errors.New("Query Error: invalid result info")
}
