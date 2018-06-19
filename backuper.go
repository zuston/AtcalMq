package main

import (
	"log"
	"flag"
	"strings"
)

const BACKUPER_DEFAULT_PATH = "/opt/aneBackup"

var (
	backuperPath = flag.String("path",BACKUPER_DEFAULT_PATH,"choose the backuper path")
	// 日期之间以 "," 来进行分割
	filterSuffix = flag.String("filter","","choose the filter backuper path suffix")
)


func init(){
	flag.Parse()
}

func main(){
	filterFiles := filter(*backuperPath,strings.Split(*filterSuffix,","))
	log.Println("current filter list is : [%s]",filterFiles)

	goroutineReader(filterFiles)
	select {

	}
}


func goroutineReader(files []string) {

}

func filter(path string, filterSuffixs []string) []string {
	return nil
}

type backuperStruct struct {
	queueName string
	info string
}


