package main

import (
	"encoding/json"
	"fmt"
)

type Server struct {
	ServerName string
	ServerIP   string
}

type Serverslice struct {
	Servers []Server
}

type CenterLoad struct{
	DataType int
	EwbNo string
	Weight float32
}

type ClList struct{

}

func main() {
	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)
	fmt.Println(s)


	strring := `[{"dataType":2,"ewbNo":"300156177594","weight":27.5}]`

	var cl []CenterLoad

	json.Unmarshal([]byte(strring), &cl)
	fmt.Println(cl[0])
}