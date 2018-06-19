package main

import (
	"log"
	"flag"
	"strings"
	"github.com/zuston/AtcalMq/util"
	"sync"
	"os"
	"bufio"
	"fmt"
	"time"
)

//const BACKUPER_DEFAULT_PATH = "/opt/aneBackup"
const BACKUPER_DEFAULT_PATH = "/temp/"

// 备份通道
var infoChannel chan backuperStruct


var (
	backuperPath = flag.String("path",BACKUPER_DEFAULT_PATH,"choose the backuper path")
	// 日期之间以 "," 来进行分割
	filterSuffix = flag.String("filter","","choose the filter backuper path suffix")
)

var barrier sync.WaitGroup


func init(){
	flag.Parse()
	infoChannel = make(chan backuperStruct,1000)
}

func main(){
	//testInc()
	//
	//return


	filterFiles := filter(*backuperPath,strings.Split(*filterSuffix,","))
	log.Printf("current filter list is : [%s]",filterFiles)

	goroutineReader(filterFiles)
	select {

	}
}


func goroutineReader(files []string) {
	for i:=0;i<len(files);i++{
		barrier.Add(1)
	}
	for _, filePath := range files{
		go func() {
			queueName,YMD := parse(filePath)
			startTime := time.Now()
			_ = YMD
			file, err := os.Open(filePath)
			// 恢复机制，记录点
			util.CheckPanic(err)

			defer func() {
				file.Close()
				barrier.Done()
			}()

			reader := bufio.NewReader(file)
			for {
				line, err := readLine(reader)
				if err!=nil {
					break
				}
				infoTag := line[0:22]
				info := line[22:]
				fmt.Println(infoTag,info)
				infoChannel <- backuperStruct{
					queueName:queueName,
					info:info,
				}
			}
			endTime := time.Now()
			fmt.Printf("%s reading cost:%s\n",queueName,endTime.Sub(startTime))
		}()
	}
	log.Println("waiting reading the files")
	barrier.Wait()

}


func parse(filePath string) (string, string) {
	arrs := strings.Split(filePath,".")
	return arrs[0],arrs[2]
}


func filter(path string, filterSuffixs []string) []string {
	var files []string
	for _,filterSuffix := range filterSuffixs{
		f := util.WalkDir(path,filterSuffix)
		for _,key := range f{
			files = append(files,key)
		}
	}
	return files
}

type backuperStruct struct {
	queueName string
	info string
}

func readLine(r *bufio.Reader) (string, error) {
	line, isprefix, err := r.ReadLine()
	for isprefix && err == nil {
		var bs []byte
		bs, isprefix, err = r.ReadLine()
		line = append(line, bs...)
	}
	return string(line), err
}


// 测试程序
func testInc(){
	files := filter("/temp",[]string{"2018-06-19"})
	log.Println(files)

}

