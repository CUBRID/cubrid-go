// Copyright 2019 The Go-CUBRID-Driver Authors. All rights reserved.
//
// Package cubrid provides a CUBRID driver for Go's database/sql package.
//
// The driver should be used via the database/sql package:
//
//  import "database/sql"
//  import _ "github.com/CUBRID/cubrid-go"

package cubrid

type col_type	byte

const (
	type_short  col_type = iota
	type_long
	type_float
	type_double
	type_bigint
	type_date
	type_time
	type_datetime
	type_timestamp
	type_varchar
	type_bit
	type_set
	type_blob
	type_clob
)

type cub_col struct {
	table		string
	col_type	byte
	precision	int
	scale		int
}

