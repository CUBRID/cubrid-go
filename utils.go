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
	"database/sql/driver"
	"time"
)

// NullTime represents a time.Time that may be NULL.
// NullTime implements the Scanner interface so
//
// This NullTime implementation is not driver-specific
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

type NullByte struct {
	data  []byte
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
// The value type must be time.Time or string
func (nt *NullTime) Scan(value interface{}) (err error) {
        if value == nil {
                nt.Time, nt.Valid = time.Now(), false
                return nil
        }
        nt.Valid = true
	switch v := value.(type) {
		case time.Time:
			nt.Time = v
		default:
			nt.Time = time.Now()
	}

        return nil
}

// NullTime Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// NullByte represents a []byte that may be NULL.
// NullByte implements the Scanner interface so
//
// This NullByte implementation is not driver-specific
func (nb *NullByte) Scan(value interface{}) (err error) {
        if value == nil {
                nb.data, nb.Valid = nil, false
                return nil
        }
        nb.Valid = true
	switch v := value.(type) {
		case []byte:
			nb.data = v
		default:
			nb.data = nil
	}

        return nil
}

// NullByte Value implements the driver Valuer interface.
func (nb NullByte) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.data, nil
}

