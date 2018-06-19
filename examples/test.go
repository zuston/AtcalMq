package main

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine

	"fmt"
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/core/pullcore"
)

func main() {
	a := "\\"
	fmt.Println(a!="\\")
	return
	line := `[{"dataType":2,"ewbNo":"300191790928","ewbsListNo":"212718313218061701","hewbNo":"30019179092800140010","nextSiteCode":"0287123","nextSiteId":18313,"operatorCode":"028656,081370","platformCode":"199","scanMan":"028656","scanTime":"2018-06-17 20:56:29","scanType":1,"siteCode":"0282000","siteId":2127,"volume":0.07,"weight":21.79},{"dataType":2,"ewbNo":"300191790928","ewbsListNo":"212718313218061701","hewbNo":"30019179092800140002","nextSiteCode":"0287123","nextSiteId":18313,"operatorCode":"028656,081370","platformCode":"199","scanMan":"028656","scanTime":"2018-06-17 20:56:29","scanType":1,"siteCode":"0282000","siteId":2127,"volume":0.07,"weight":21.79},{"dataType":2,"ewbNo":"300191790928","ewbsListNo":"212718313218061701","hewbNo":"30019179092800140006","nextSiteCode":"0287123","nextSiteId":18313,"operatorCode":"028656,081370","platformCode":"199","scanMan":"028656","scanTime":"2018-06-17 20:56:30","scanType":1,"siteCode":"0282000","siteId":2127,"volume":0.07,"weight":21.79},{"dataType":2,"ewbNo":"300191790928","ewbsListNo":"212718313218061701","hewbNo":"30019179092800140007","nextSiteCode":"0287123","nextSiteId":18313,"operatorCode":"028656,081370","platformCode":"199","scanMan":"028656","scanTime":"2018-06-17 20:56:32","scanType":1,"siteCode":"0282000","siteId":2127,"volume":0.07,"weight":21.79},{"dataType":2,"ewbNo":"300191790928","ewbsListNo":"212718313218061701","hewbNo":"30019179092800140012","nextSiteCode":"0287123","nextSiteId":18313,"operatorCode":"028656,081370","platformCode":"199","scanMan":"028656","scanTime":"2018-06-17 20:56:32","scanType":1,"siteCode":"0282000","siteId":2127,"volume":0.07,"weight":21.79},{"dataType":2,"ewbNo":"300191790928","ewbsListNo":"212718313218061701","hewbNo":"30019179092800140001","nextSiteCode":"0287123","nextSiteId":18313,"operatorCode":"028656,081370","platformCode":"199","scanMan":"028656","scanTime":"2018-06-17 20:56:33","scanType":1,"siteCode":"0282000","siteId":2127,"volume":0.07,"weight":21.79},{"dataType":2,"ewbNo":"300191790928","ewbsListNo":"212718313218061701","hewbNo":"30019179092800140011","nextSiteCode":"0287123","nextSiteId":18313,"operatorCode":"028656,081370","platformCode":"199","scanMan":"028656","scanTime":"2018-06-17 20:56:34","scanType":1,"siteCode":"0282000","siteId":2127,"volume":0.07,"weight":21.79}]`
	mappers := (pullcore.ModelGen([]byte(line)))
	for _, mapper := range mappers{
		cmapper := pullcore.V2Byte(mapper)
		v, ok := cmapper["vehicleno"]
		fmt.Println(v,ok)
	}
	return

	fmt.Println(util.IsDir("/Users/zuston/goDev/src/github.com/zuston/AtcalMq"))
	fmt.Println(util.WalkDir("/Users/zuston/goDev/src/github.com/zuston/AtcalMq",".model"))
	return
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