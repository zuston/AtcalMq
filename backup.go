package main

import (
	"github.com/zuston/AtcalMq/util"
	"fmt"
	"os"
	"bufio"
	"time"
	"strings"
	"github.com/zuston/AtcalMq/core/pullcore"
)

type sendObj struct {
	queueName string
	info string
}

var channel chan sendObj
var done chan bool
var finalDone chan bool

var statics map[string]int64

func init(){
	channel = make(chan sendObj,1000)
	done = make(chan bool,1)
	finalDone = make(chan bool,1)

	statics = make(map[string]int64,10)
}

// 备份数据补足入库代码

func main() {
	filterFilesList := filterFiles("/tmp/backup","backup_2018_6_4","backup_2018_6_3")

	fmt.Println(filterFilesList)

	go backUpHandler()

	goroutineRead2Hbase(filterFilesList)

	select {

	}
}


func goroutineRead2Hbase(filePaths []string) {
		for _, filePath := range filePaths{

			func() {
				queueName := splitPath(filePath)
				startTime := time.Now()
				f, err := os.Open(filePath)
				util.CheckPanic(err)
				defer f.Close()
				r := bufio.NewReader(f)
				for {
					line, err := readLine(r)
					if err!=nil {
						break
					}
					InfoTag := line[0:28]
					if !strings.Contains(InfoTag,"2018/06/03") && !strings.Contains(InfoTag,"2018/06/04") {
						continue
					}
					line = line[28:]
					channel <- sendObj{
						queueName:queueName,
						info:line,
					}
				}
				endTime := time.Now()
				fmt.Printf("%s 读取耗时间:%s\n",queueName,endTime.Sub(startTime))

			}()

		}
		//time.Sleep(time.Second*2)
		done <- false
}


func splitPath(s string) string {
	index_1 := strings.LastIndex(s,"/")
	index_2 := strings.Index(s,".")
	return s[index_1+1:index_2]
}

func filterFiles(path string, filterSuffixs... string) []string{
	var files []string
	for _,filterSuffix := range filterSuffixs{
		f := util.WalkDir(path,filterSuffix)
		for _,key := range f{
			files = append(files,key)
		}
	}
	return files
}

func ReadLine(filename string) {
	f, _ := os.Open(filename)
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		_, err := readLine(r)
		if err != nil {
			break
		}
	}

}

// 解决一行过大
func readLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}


func backUpHandler(){
	startTime := time.Now()

	go func() {
		<- done
		close(channel)
	}()

	for sendInfo := range channel{
		queueName := sendInfo.queueName
		info := sendInfo.info

		list := pullcore.ModelGen([]byte(info))
		value,ok := statics[queueName]
		if !ok {
			statics[queueName] = 1
		}else{
			statics[queueName] = value+int64(len(list))
		}
	}
	endTime := time.Now()
	fmt.Printf("队列总计处理时间%s",endTime.Sub(startTime))
	fmt.Println()
	fmt.Println(statics)
}