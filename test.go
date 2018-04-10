package main

import (
	"encoding/json"
	"fmt"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"github.com/zuston/AtcalMq/core"
	"os"
	"bufio"
	"strings"
	"io"
	"time"
	"github.com/zuston/AtcalMq/rabbitmq"
	"github.com/json-iterator/go"
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

	sd := core.UidGen()
	fmt.Println(sd)
	return

	str := `
		{"ServerName":{"ip":12,"op":"iou"},"ServerIP":"123"}
	`
	v := &Server{}
	jsoniter.Unmarshal([]byte(str),&v)
	fmt.Println(v.ServerIP)
	return

	go func() {
		time.Sleep(5*time.Second)
		fmt.Println("output go")
	}()
	defer func() {
		fmt.Println("defer")
	}()

	return

	go rabbitmq.NewWatcher()
	select {

	}
	mapper := make(map[string]string,4)
	ffile, err := os.Open("/opt/mq.cfg")

	buffer := bufio.NewReader(ffile)
	for {
		line, err := buffer.ReadString('\n')
		if err!= nil {
			if err==io.EOF {
				fmt.Println("结束")
				return
			}
			return
		}
		line = strings.TrimSpace(line)
		fmt.Println(line,"*")
		tagIndex := strings.Index(line,"=")
		if tagIndex==-1 {
			fmt.Println("error, tagIndex=",tagIndex)
			return
		}
		key := strings.TrimSpace(line[:tagIndex])
		value := strings.TrimSpace(line[tagIndex+1:])
		mapper[key] = value

	}
	return

	values := map[string]map[string][]byte{"base": map[string][]byte{"ScanMan": []byte{0}}}
	putRequest, _ := hrpc.NewPutStr(context.Background(), "centerload","1", values)
	_, err = core.Hconn.Client.Put(putRequest)
	if err!=nil{
		fmt.Println(err)
	}
	return

	var s Serverslice
	p := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(p), &s)
	fmt.Println(s)


	strring := `[{"dataType":2,"ewbNo":"300156177594","weight":27.5}]`

	var cl []CenterLoad

	json.Unmarshal([]byte(strring), &cl)
	fmt.Println(cl[0].DataType)
}