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
	"unsafe"
	"strconv"
	"time"
	"encoding/hex"
)

type cub_result struct {
        columns []cub_col
        names   []string
        affected_rows   int64
        last_insert     int64
}

type cub_rows struct {
	conn	int
	handle	int
	result	*cub_result
}

func (row *cub_rows) Columns() []string {
	if row.result.names != nil {
		return row.result.names
	}

	return nil
}

func (r *cub_rows) Close() error {
	return nil
}

func (r *cub_rows) Next(dest []driver.Value) error {
	var err_buf	C.T_CCI_ERROR
	var indicator	C.int
	var value*	byte
	var blob_data	[1024]C.char
	var blob	C.T_CCI_BLOB
	var v		[]byte = nil
	var err		error = nil
	var blob_start	C.longlong = 0

	res := C.cci_cursor(C.int(r.handle), 1, 1, &err_buf)
	if (res < 0) {
		return errors.New(C.GoString(&err_buf.err_msg[0]))
	}

	res = C.cci_fetch(C.int(r.handle), &err_buf)
	if (res < 0) {
		return errors.New(C.GoString(&err_buf.err_msg[0]))
	}

	for i := range dest {
		if r.result.columns[i].col_type == C.CCI_U_TYPE_BLOB {
			res = C.cci_get_data(C.int(r.handle), C.int(i + 1), C.CCI_A_TYPE_BLOB, unsafe.Pointer(&blob), &indicator)
		} else {
			res = C.cci_get_data(C.int(r.handle), C.int(i + 1), C.CCI_A_TYPE_STR, unsafe.Pointer(&value), &indicator)
			if indicator > 0 {
				v = C.GoBytes(unsafe.Pointer(value), indicator)
			}
		}
		switch r.result.columns[i].col_type {
			case C.CCI_U_TYPE_BIT, C.CCI_U_TYPE_VARBIT:
				dst := make([]byte, hex.DecodedLen(len(v)))
				dst, err = hex.DecodeString(string(v))
				dest[i] = dst
			case C.CCI_U_TYPE_INT:
				dest[i], err = strconv.Atoi(string(v))
			case C.CCI_U_TYPE_SHORT:
				dest[i], err = strconv.Atoi(string(v))
			case C.CCI_U_TYPE_BIGINT:
				dest[i], err = strconv.ParseInt(string(v), 10, 64)
			case C.CCI_U_TYPE_FLOAT:
				dest[i], err = strconv.ParseFloat(string(v), 32)
			case C.CCI_U_TYPE_DOUBLE:
				dest[i], err = strconv.ParseFloat(string(v), 64)
			case C.CCI_U_TYPE_TIME:
				dest[i], err = time.Parse("15:04:05", string(v))
			case C.CCI_U_TYPE_DATE:
				dest[i], err = time.Parse("2006-01-02", string(v))
			case C.CCI_U_TYPE_DATETIME:
				dest[i], err = time.Parse("2006-01-02 15:04:05.000", string(v))
			case C.CCI_U_TYPE_TIMESTAMP:
				dest[i], err = time.Parse("2006-01-02 15:04:05", string(v))
			case C.CCI_U_TYPE_BLOB:
				n := C.cci_blob_size(blob)
				res = C.cci_blob_read(C.int(r.conn), blob, blob_start, C.int(n), &blob_data[0], &err_buf)
				if res < 0 {
					err = errors.New(C.GoString(&err_buf.err_msg[0]))
				} else {
					dest[i] = C.GoBytes(unsafe.Pointer(&blob_data[0]), C.int(n))
				}
			case C.CCI_U_TYPE_CLOB:
				dest[i] = nil
			default:
				dest[i] = string(v)
		}

		//fmt.Println(dest[i])
	}
	return err
}

func (r *cub_result) LastInsertId() (int64, error) {
        return r.last_insert, nil
}

func (r *cub_result) RowsAffected() (int64, error) {
	return r.affected_rows, nil
}

