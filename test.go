package main

import (
	"encoding/json"
	"fmt"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"github.com/zuston/AtcalMq/core"
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

	values := map[string]map[string][]byte{"base": map[string][]byte{"ScanMan": []byte{0}}}
	putRequest, _ := hrpc.NewPutStr(context.Background(), "centerload","1", values)
	_, err := core.Hconn.Client.Put(putRequest)
	if err!=nil{
		fmt.Println(err)
	}
	return

	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)
	fmt.Println(s)


	strring := `[{"dataType":2,"ewbNo":"300156177594","weight":27.5}]`

	var cl []CenterLoad

	json.Unmarshal([]byte(strring), &cl)
	fmt.Println(cl[0].DataType)
}