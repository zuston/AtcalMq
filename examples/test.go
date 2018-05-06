package main

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine

	"fmt"
)

func main() {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "zuston", "todo")
	db.Register("set names utf8")

	err := db.Connect()
	if err != nil {
		panic(err)
	}
	rows, res, err := db.Query("select * from user")
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		fmt.Println(res.Map("id"))
		for _, col := range row {
			if col == nil {
				// col has NULL value
			} else {
				// Do something with text in col (type []byte)
			}
		}

		//fmt.Println(len(row))
		//// You can get specific value from a row
		//val1 := row[1].([]byte)
		//
		//// You can use it directly if conversion isn't needed
		//os.Stdout.Write(val1)
		//
		//// You can get converted value
		//number := row.Int(0)      // Zero value
		//str    := row.Str(1)      // First value
		//bignum := row.MustUint(2) // Second value
		//
		//// You may get values by column name
		//first := res.Map("FirstColumn")
		//second := res.Map("SecondColumn")
		//val1, val2 := row.Int(first), row.Str(second)
	}
}