package main

import (
	ui "github.com/gizak/termui"
	"net/rpc"
	"github.com/zuston/AtcalMq/util"
	"time"
	"github.com/zuston/AtcalMq/rabbitmq"
	"encoding/json"
	"strconv"
	"fmt"
	"log"
)

const RPC_URL = "127.0.0.1:9898"

var(
	clog *util.Logger
	rpcClient *rpc.Client
	statusRpc string
)

func init(){
	var err error
	rpcClient, err = rpc.DialHTTP("tcp",RPC_URL)
	util.CheckPanic(err)
}

func main(){
	go timelyGetStatus()
	uiInit()
}

func timelyGetStatus(){
	timeTicker := time.NewTicker(time.Second*5)
	var rpcRes string
	for _ = range timeTicker.C{
		err := rpcClient.Call("Watcher.GetAll","", &rpcRes)
		if !util.CheckErr(err){
			continue
		}
		statusRpc = rpcRes
	}
}

func uiInit(){

	pullCD,_ := initCoverData()

	if err := ui.Init(); err != nil {
		panic(err)
	}
	defer ui.Close()


	consoleData := [][]string{
		[]string{
			"time:2018-23-83",
			"pushQueueNumber:434",
			"pullQueueNumber:0",
		},
	}

	// daemon process status
	console := ui.NewTable()
	console.Rows = consoleData
	console.FgColor = ui.ColorBlack
	console.BgColor = ui.ColorDefault
	console.TextAlign = ui.AlignCenter
	console.Separator = true
	console.Analysis()
	console.SetSize()
	console.BgColors[0] = ui.ColorGreen
	console.Border = false
	ui.Render(console)

	coverData := [][]string{
		[]string{
			"QUEUE_NAME",
			"STOCK",
			"UNIT_HANDLE",
			"TOTAL",
		},
	}

	pull := ui.NewTable()
	pull.Rows = pullCD
	pull.FgColor = ui.ColorWhite
	pull.BgColor = ui.ColorDefault
	pull.TextAlign = ui.AlignLeft
	pull.Separator = false
	pull.Analysis()
	pull.SetSize()
	pull.Border = true
	pull.BorderLabel = "consumer queue status"
	ui.Render(pull)

	push := ui.NewTable()
	push.Rows = coverData
	push.FgColor = ui.ColorWhite
	push.BgColor = ui.ColorDefault
	push.TextAlign = ui.AlignCenter
	push.Separator = false
	push.Analysis()
	push.SetSize()
	push.BorderLabel="producer queue status"
	push.Border = true

	ui.Render(push)

	// build layout
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, console),
		),
		ui.NewRow(
				ui.NewCol(6, 0, pull),
				ui.NewCol(6, 0, push),
			))

	// calculate layout
	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})
	ui.Handle("/sys/wnd/resize", func(e ui.Event) {
		ui.Body.Width = ui.TermWidth()
		ui.Body.Align()
		ui.Clear()
		ui.Render(ui.Body)
	})

	// 处理刷新逻辑
	ui.Handle("/timer/1s", func(e ui.Event) {
		var rpcObj rabbitmq.RpcObj
		json.Unmarshal([]byte(statusRpc),&rpcObj)

		var pullObjs [][]string
		pullHeader := []string{
			"QUEUE_NAME",
			"STOCK",
			"UNIT_HANDLE",
			"TOTAL",
		}
		pullObjs = append(pullObjs,pullHeader)
		pull.BgColors[0] = ui.ColorRed
		for _,pushQueueObj := range rpcObj.PullSupervisorObjs{
			var pushObj = []string{
				pushQueueObj.QueueName,
				pushQueueObj.Overstock,
				strconv.Itoa(pushQueueObj.UnitHandlerAblitity),
				fmt.Sprintf("%d",pushQueueObj.TotalHandlerCount),
			}
			pullObjs = append(pullObjs,pushObj)
		}
		pull.Rows = pullObjs
		ui.Render(pull)

		var consoleObjs [][]string
		runtime := rpcObj.Runtime
		consoleObj := []string{
			fmt.Sprintf("runTime:%d小时%d分钟%d秒",int(runtime.Hours()),int(runtime.Minutes()),int(runtime.Seconds()/60)),
			fmt.Sprintf("push queue number:%d",rpcObj.PushQueueNumber),
			fmt.Sprintf("pull queue number:%d",rpcObj.PullQueueNumber),
		}
		consoleObjs = append(consoleObjs,consoleObj)
		console.Rows = consoleObjs
		ui.Render(console)
	})

	ui.Loop()
}

// 初始化填充数据
func initCoverData() ([][]string,[][]string){
	log.Println("init the cover data.....")
	for{
		if statusRpc!="" {
			var rpcObj rabbitmq.RpcObj
			json.Unmarshal([]byte(statusRpc),&rpcObj)
			var cover_1 [][]string
			for i:=0;i<len(rpcObj.PullSupervisorObjs)+1;i++{
				data := []string{
					"","","","",
				}
				cover_1 = append(cover_1,data)
			}
			return cover_1,nil
		}
		log.Println("loading.....")
		time.Sleep(5*time.Second)
	}
}


func TestRpc(){
	fmt.Println(statusRpc)
	var rpcObj rabbitmq.RpcObj
	json.Unmarshal([]byte(statusRpc),&rpcObj)
	fmt.Println(rpcObj.PushQueueNumber)
	var pushObjs [][]string
	for _,pushQueueObj := range rpcObj.PullSupervisorObjs{
		var pushObj = []string{
			pushQueueObj.QueueName,
			pushQueueObj.Overstock,
			strconv.Itoa(pushQueueObj.UnitHandlerAblitity),
		}
		pushObjs = append(pushObjs,pushObj)
	}
	fmt.Println(pushObjs)
}