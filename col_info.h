// Copyright 2019 The Go-CUBRID-Driver Authors. All rights reserved.
//
// Package cubrid provides a CUBRID driver for Go's database/sql package.
//
// The driver should be used via the database/sql package:
//
//  import "database/sql"
//  import _ "github.com/CUBRID/cubrid-go"

// This struct type is a clone of T_CCI_COL_INFO
// the member ext_type is replacement of 'type'
// in T_CCI_COL_INFO for avoiding go's reserved word
typedef struct
  {
    T_CCI_U_TYPE ext_type;
    char is_non_null;
    short scale;
    int precision;
    char *col_name;
    char *real_attr;
    char *class_name;
    char *default_value;
    char is_auto_increment;
    char is_unique_key;
    char is_primary_key;
    char is_foreign_key;
    char is_reverse_index;
    char is_reverse_unique;
    char is_shared;
  } T_CCI_COL9x_INFO;

typedef struct
  {
    unsigned char ext_type;
    char is_non_null;
    short scale;
    int precision;
    char *col_name;
    char *real_attr;
    char *class_name;
    char *default_value;
    char is_auto_increment;
    char is_unique_key;
    char is_primary_key;
    char is_foreign_key;
    char is_reverse_index;
    char is_reverse_unique;
    char is_shared;
  } T_CCI_COL10_INFO;

