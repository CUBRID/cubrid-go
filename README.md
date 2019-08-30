# cubrid-go
CUBRID Go Driver on top of CCI
> Export following environment variables to build go applications
```bash
export CGO_CFLAGS="-I$CUBRID/include"
export CGO_LDFLAGS="-L$CUBRID/lib -lcascci -lnsl -lpthread -lrt"

go get -u github.com/CUBRID/cubrid-go
```
The following code is an example application using CUBRID go driver.

<pre>
<code>
package main

import (
    "database/sql"
    _ "github.com/CUBRID/cubrid-go"
    "log"
    "fmt"
    "time"
)

func main() {
        db, err := sql.Open("cubrid", "cci:CUBRID:localhost:55300:demodb:dba::")
        if err != nil {
                log.Fatal(err)
        }
        defer db.Close()

        var id  int
        var a_bit       []byte
        var b_vbit      []byte
        var c_num       float64
        var d_float     float32
        var e_double    float64
        var f_date      time.Time
        var g_time      time.Time
        var g_timest    time.Time
        var h_set       string
        var i_bigint    int64
        var j_datetm    time.Time
        var k_blob      sql.NullString
        var l_clob      sql.NullString

        rows, err := db.Query("select * from tbl_go")

        if err != nil {
                log.Fatal(err)
        }

        defer rows.Close()

        for rows.Next() {
                err := rows.Scan(&id, &a_bit, &b_vbit, &c_num, &d_float, &e_double,
                                 &f_date, &g_time, &g_timest, &h_set, &i_bigint, &j_datetm, &k_blob, &l_clob)
                if err != nil {
                        log.Fatal(err)
                }

                fmt.Println(id)
                fmt.Println(a_bit)
                fmt.Println(b_vbit)
                fmt.Println(c_num)
                fmt.Println(d_float)
                fmt.Println(e_double)
                fmt.Println(f_date)
                fmt.Println(g_time)
                fmt.Println(g_timest)
                fmt.Println(h_set)
                fmt.Println(i_bigint)
                fmt.Println(j_datetm)
                fmt.Println(k_blob)
                fmt.Println(l_clob)
        }
}
</code>
</pre>
