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
	"strconv"
	"time"
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

func (s *cub_stmt) bind(args []driver.Value) int {
	var handle	C.int
	var parmNum	C.int
	var pLonglong	C.longlong
	var pDouble	C.double
	var pString	*C.char
	var blob	C.T_CCI_BLOB
	var err		C.T_CCI_ERROR

	handle = C.int(s.handle)
	for i, arg := range args {
		parmNum = C.int(i) + 1
		switch v := arg.(type) {
			case int64:
				pLonglong = C.longlong(v)
				C.cci_bind_param (handle, parmNum, C.CCI_A_TYPE_BIGINT,
					unsafe.Pointer(&pLonglong), C.CCI_U_TYPE_BIGINT, C.CCI_BIND_PTR);
			case float64:
				pDouble = C.double(v)
				C.cci_bind_param (handle, parmNum, C.CCI_A_TYPE_DOUBLE,
					unsafe.Pointer(&pDouble), C.CCI_U_TYPE_DOUBLE, C.CCI_BIND_PTR);
			case string:
				pString = C.CString(v)
				C.cci_bind_param (handle, parmNum, C.CCI_A_TYPE_STR,
					unsafe.Pointer(pString), C.CCI_U_TYPE_STRING, C.CCI_BIND_PTR);
			case []byte:
				if v != nil {
					size := len (v)
					C.cci_blob_new(C.int(s.conn.handle), &blob, &err)
					C.cci_blob_write(C.int(s.conn.handle), blob, 0, C.int(size), (*C.char) (unsafe.Pointer(&v[0])), &err)
					C.cci_bind_param (handle, parmNum, C.CCI_A_TYPE_BLOB,
						unsafe.Pointer(blob), C.CCI_U_TYPE_BLOB, C.CCI_BIND_PTR);
				} else {
					C.cci_bind_param (handle, parmNum, C.CCI_A_TYPE_BLOB,
						unsafe.Pointer(nil), C.CCI_U_TYPE_BLOB, C.CCI_BIND_PTR);
				}
			case time.Time:
				pString = C.CString(v.Format("2006-01-02 15:04:05.000"))
				C.cci_bind_param (handle, parmNum, C.CCI_A_TYPE_STR,
					unsafe.Pointer(pString), C.CCI_U_TYPE_STRING, C.CCI_BIND_PTR);
			case nil:
				C.cci_bind_param (handle, parmNum, C.CCI_A_TYPE_STR,
					unsafe.Pointer(nil), C.CCI_U_TYPE_STRING, C.CCI_BIND_PTR);
			default:
				break
		}
	}

	return 0
}

func (s *cub_stmt) Exec(args []driver.Value) (driver.Result, error) {
	var handle	C.int
	var flag	C.char
	var err_buf	C.T_CCI_ERROR
	var last_insert_id_ptr *C.char
	var last_insert_id int = 0

	s.params = len(args)

	if len(args) > 0 {
		if (s.bind(args) < 0) {
			err := errors.New("Exec: some parameter cannot be converted to vaild DB type")
			return nil, err
		}
	}

	handle = C.int(s.handle)
	flag = C.char(s.flag)

	res := C.cci_execute(handle, flag, 0, &err_buf)
	if int(res) < 0 {
		err := errors.New(C.GoString(&err_buf.err_msg[0]))
		return nil, err
	}

	cci_ret := C.cci_get_last_insert_id(C.int(s.conn.handle), unsafe.Pointer(&last_insert_id_ptr), &err_buf)

	if C.int(cci_ret) == 0 {
		last_insert_id, _ = strconv.Atoi(C.GoString(last_insert_id_ptr))
	}

	return &cub_result {
		affected_rows: int64(res),
		last_insert: int64(last_insert_id),
	}, nil
}

func (s *cub_stmt) Query(args []driver.Value) (driver.Rows, error) {
	var handle	C.int
	var flag	C.char
	var err_buf	C.T_CCI_ERROR
	var col_nums	C.int
	var col_info*	C.T_CCI_COL_INFO
	var ci9x	[]C.T_CCI_COL9x_INFO = nil
	var ci10	[]C.T_CCI_COL10_INFO = nil
	var stmt_type	C.T_CCI_CUBRID_STMT

	s.params = len(args)

	if len(args) > 0 {
		if s.bind(args) < 0 {
			err := errors.New("Query: some parameter cannot be converted to vaild DB type")
			return nil, err
		}
	}

	handle = C.int(s.handle)
	flag = C.char(s.flag)
	col_info = C.cci_get_result_info(handle, &stmt_type, &col_nums)

	length := int(col_nums)
	hdr := reflect.SliceHeader {
		Data:	uintptr(unsafe.Pointer(col_info)),
		Len:	length,
		Cap:	length,
	}

	if C.sizeof_T_CCI_COL_INFO == C.sizeof_T_CCI_COL9x_INFO {
		ci9x = *(*[]C.T_CCI_COL9x_INFO)(unsafe.Pointer(&hdr))
	} else {
		ci10 = *(*[]C.T_CCI_COL10_INFO)(unsafe.Pointer(&hdr))
	}

	if ci9x != nil || ci10 != nil {
		col_list := make([]cub_col, col_nums)
		name_list := make([]string, col_nums)
		if ci9x != nil {
			for i:= 0; i < int(col_nums); i++ {
				name_list[i] = C.GoString(ci9x[i].col_name)
				col_list[i].table = C.GoString(ci9x[i].class_name)
				col_list[i].col_type = byte(ci9x[i].ext_type)
				col_list[i].scale = int(ci9x[i].scale)
				col_list[i].precision = int(ci9x[i].precision)
			}
		} else {
			for i:= 0; i < int(col_nums); i++ {
				name_list[i] = C.GoString(ci10[i].col_name)
				col_list[i].table = C.GoString(ci10[i].class_name)
				col_list[i].col_type = byte(ci10[i].ext_type)
				col_list[i].scale = int(ci10[i].scale)
				col_list[i].precision = int(ci10[i].precision)
			}
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
