# cubrid-go
CUBRID Go Driver on top of CCI

### Requirements
- git
- golang
- CUBRID CCI Driver (any version)

### Install required packages
```
$ sudo yum install git
$ sudo yum install go
```

### Install CUBRID CCI Driver
- Install CUBRID CCI Driver or CUBRID Engine
- download from `http://ftp.cubrid.org`

### Exports Environment Variables for Go Build
- CUBRID Engine 11.2 or higher (or CCI 11.0.0)

```
export CGO_CFLAGS="-I$CUBRID/cci/include"
export CGO_LDFLAGS="-L$CUBRID/cci/lib -lcascci -lnsl -lpthread -lrt"
export LD_LIBRARY_PATH=$CUBRID/cci/lib:$LD_LIBRARY_PATH
```

- Other versions (versions less than 11.2 or CCI 11.0.0)

```
export CGO_CFLAGS="-I$CUBRID/cci/include"
export CGO_LDFLAGS="-L$CUBRID/cci/lib -lcascci -lnsl -lpthread -lrt"
export LD_LIBRARY_PATH=$CUBRID/lib:$LD_LIBRARY_PATH
```

### Downloads and install CUBRID Go Driver

```
$ go env -w GO111MODULE=off
$ go get -u github.com/CUBRID/cubrid-go
```

### Run CUBRID Engine & demodb with following parameters for testing

- IP address: 127.0.0.1
- Broker Port: 33000

```
$ cubrid service start
$ cd $CUBRID/demo
$ sh make_cubrid_demo.sh
$ cubrid server start demodb
```

### Compile and run sample application with CUBRID Go Driver

- drag and drop following go sample code as olympic.go

```
$ vi olympic.go
$ go build olympic.go
$ ./olympic
```

<pre>
<code>
package main

import (
	"database/sql"
	"github.com/CUBRID/cubrid-go"
	"log"
	"fmt"
)

func main() {
	db, err := sql.Open("cubrid", "cci:CUBRID:127.0.0.1:33000:demodb:public::")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var year		sql.NullInt32
	var city		sql.NullString
	var open_date		cubrid.NullTime

	rows, err := db.Query("select host_year, host_city, opening_date from olympic")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&year, &city, &open_date)
		if err != nil {
		log.Fatal(err)
		}

		if year.Valid {
			fmt.Print(year.Int32, "\t")
		}

		if city.Valid {
			fmt.Print(city.String, "\t")
			if  len(city.String) < 8 {
				fmt.Print("\t")
			}
		}

		if open_date.Valid {
			fmt.Println(open_date.Time)
		}
	}
}

</code>
</pre>

 - Expected Result
<pre>
<code>
2004	Athens		2004-08-13 00:00:00 +0000 UTC
2000	Sydney		2000-09-15 00:00:00 +0000 UTC
1996	Atlanta		1996-07-19 00:00:00 +0000 UTC
1992	Barcelona	1992-07-25 00:00:00 +0000 UTC
1988	Seoul		1988-09-17 00:00:00 +0000 UTC
1984	Los Angeles	1984-07-28 00:00:00 +0000 UTC
1980	Moscow		1980-07-19 00:00:00 +0000 UTC
1976	Montreal	1976-07-17 00:00:00 +0000 UTC
1972	Munich		1972-08-26 00:00:00 +0000 UTC
1968	Mexico City	1968-10-12 00:00:00 +0000 UTC
1964	Tokyo		1964-10-10 00:00:00 +0000 UTC
1960	Rome		1960-08-25 00:00:00 +0000 UTC
1956	Melbourne	1956-11-22 00:00:00 +0000 UTC
1952	Helsinki	1952-07-19 00:00:00 +0000 UTC
1948	London		1948-07-29 00:00:00 +0000 UTC
1936	Berlin		1936-08-01 00:00:00 +0000 UTC
1932	Los Angeles	1932-07-30 00:00:00 +0000 UTC
1928	Amsterdam	1928-07-28 00:00:00 +0000 UTC
1924	Paris		1924-05-04 00:00:00 +0000 UTC
1920	Antwerp		1920-04-20 00:00:00 +0000 UTC
1912	Stockholm	1912-05-05 00:00:00 +0000 UTC
1908	London		1908-04-27 00:00:00 +0000 UTC
1904	St. Louis	1904-07-01 00:00:00 +0000 UTC
1900	Paris		1900-05-14 00:00:00 +0000 UTC
1896	Athens		1896-04-06 00:00:00 +0000 UTC
</code>
</pre>

- The following code is an example application using CUBRID go driver.

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

	var id  int     sql.NullInt64
	var a_bit       cubrid.NullByte
	var b_vbit      cubrid.NullByte
	var c_num       sql.NullFloat64
	var d_float     sql.NullFloat64
	var e_double    sql.NullFloat64
	var f_date      cubrid.NullTime
	var g_time      cubrid.NullTime
	var g_timest    cubrid.NullTime
	var h_set       sql.NullString
	var i_bigint    sql.NullInt64
	var j_datetm    cubrid.NullTime
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
