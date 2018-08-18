package main


import (
	"net/rpc"
	"github.com/zuston/AtcalMq/util"

	"time"
	"github.com/gizak/termui"
	"encoding/json"
	"strconv"
	"fmt"
	"github.com/zuston/AtcalMq/rabbitmq"
)

var (
	zllogger *util.Logger
	client *rpc.Client

	currentInfos string
)

const CONSOLE_LOG_PATH  = "/extdata/log/AneConsole.log"


func init(){
	zllogger, _ = util.NewLogger(util.DEBUG_LEVEL,CONSOLE_LOG_PATH)
	zllogger.SetDebug()

	var err error
	client, err = rpc.DialHTTP("tcp","127.0.0.1:9898")
	if err!=nil {
		zllogger.Error("connect to the rpc server error : %s",err)
		return
	}
	zllogger.Info("connect to the rpc server success")
}

func dataCover(){
	var jsonReply string

	for{
		err := client.Call("Watcher.GetAll","", &jsonReply)
		if err!=nil {
			zllogger.Error("call function error : ",err)
			continue
		}
		//zllogger.Info("queueName : %s, info : %s",core.QUEUE_CENTERLOAD,jsonReply)
		currentInfos = jsonReply
		time.Sleep(5*time.Second)
	}
	return
}




func main() {

	go dataCover()

	if err := termui.Init(); err != nil {
		panic(err)
	}
	defer termui.Close()


	strs := []string{
		"[0] github.com/gizak/termui",
		"[1] [你好，世界](fg-blue)",
		"[2] [こんにちは世界](fg-red)",
		"[3] [color output](fg-white,bg-green)",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] nsf/termbox-go"}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = 30
	ls.Width = 70
	ls.X = 121
	ls.Y = 0
	termui.Render(ls)

	bc := termui.NewBarChart()
	bclabels := []string{"S0", "S1", "S2", "S3", "S4", "S5"}
	bc.BorderLabel = "Queue Status (PRESS q TO QUIT)"
	bc.Width = 120
	bc.Height = 30
	bc.DataLabels = bclabels
	bc.TextColor = termui.ColorGreen
	bc.BarColor = termui.ColorRed
	bc.NumColor = termui.ColorYellow

	bc.BarWidth = 5
	termui.Render(bc)


	ha := termui.NewBarChart()
	ha.BorderLabel = "Queue Unit Time Handle Ability"
	ha.Width = 120
	ha.Height = 15
	ha.DataLabels = []string{"S0"}
	ha.TextColor = termui.ColorGreen
	ha.BarColor = termui.ColorRed
	ha.NumColor = termui.ColorYellow
	ha.X = 0
	ha.Y = 31
	ha.BarWidth = 5
	termui.Render(ha)


	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	initBarLabelTag := true
	termui.Handle("/timer/1s", func(e termui.Event) {
		//t := e.Data.(termui.EvtTimer)
		var jv []rabbitmq.PullSupervisorObj
		json.Unmarshal([]byte(currentInfos),&jv)
		if initBarLabelTag {
			var labels []string
			tagIndex := 1
			var listLabels [] string

			for _,v := range jv{
				labels = append(labels,fmt.Sprintf("Q%d",tagIndex))
				listLabels = append(listLabels,fmt.Sprintf("[Q%d] %s",tagIndex,v.QueueName))
				tagIndex++
			}
			initBarLabelTag = false
			bc.DataLabels = labels
			ha.DataLabels = labels
			ls.Items = listLabels
		}
		// overstock array
		var overstockData []int
		// unit ablity arr
		var unitAbilityData []int
		for _,v := range jv{
			intv, _ := strconv.Atoi(v.Overstock)
			overstockData = append(overstockData,intv)
			unitAbilityData = append(unitAbilityData,v.UnitHandlerAblitity)
		}
		bc.Data = overstockData
		ha.Data = unitAbilityData
		termui.Render(bc)
		termui.Render(ls)
		termui.Render(ha)
	})

	termui.Loop()
}