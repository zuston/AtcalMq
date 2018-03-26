package util

import (
	"log"
	"os"
	"fmt"
	"strings"
	"time"
)

/**
support the log file splited by day
 */

const (
	DEBUG_LEVEL int = iota
	INFO_LEVEL
	WARN_LEVEL
	ERROR_LEVEL
)

type Logger struct {
	level int
	logTime int64
	fileName string
	fileFolder string
	fileHandle *os.File
	nlogger *log.Logger
	debug bool
}


func NewLogger(level int,filePath string) (*Logger,error){

	logger := &Logger{}
	logger.level = level
	logger.debug = false
	logger.nlogger = log.New(logger,"[info]",log.Ldate | log.Ltime)

	logger.fileFolder = "./"

	if index := strings.LastIndex(filePath,"/"); index!=-1 {
		logger.fileFolder = filePath[0:index] + "/"
		if index!=len(filePath) {
			logger.fileName = filePath
		}
	}

	logger.fileName = filePath

	return logger,nil
}

func (logger *Logger) SetDebug(){
	logger.debug = true
}

//rewrite the implementation of log output function
func (logger *Logger) Write(buf []byte) (n int, err error) {
	if logger.fileName=="" {
		fmt.Println("console : %s", buf)
		return  len(buf), nil
	}

	logger.handlerFile()

	if logger.fileHandle==nil {
		return len(buf),nil
	}
	if logger.debug {
		fmt.Println(string(buf))
	}
	return logger.fileHandle.Write(buf)
}

func (logger *Logger) SetLevel(level int){
	logger.level = level
}

func (logger *Logger)Debug(format string, args ...interface{}){
	if logger.level <= DEBUG_LEVEL {
		logger.nlogger.SetPrefix("Debug ")
		outLine := fmt.Sprintf(format, args...)
		if logger.debug {
			fmt.Println(outLine)
		}
		logger.nlogger.Output(2,outLine )
	}
}


func (logger *Logger)Info(format string, args ...interface{}){
	if logger.level <= INFO_LEVEL {
		logger.nlogger.SetPrefix("Info ")
		logger.nlogger.Output(2, fmt.Sprintf(format, args...))
	}
}

func (logger *Logger)Warn(format string, args ...interface{}){
	if logger.level <= WARN_LEVEL {
		logger.nlogger.SetPrefix("Warn ")
		logger.nlogger.Output(2, fmt.Sprintf(format, args...))
	}
}

func (logger *Logger)Error(format string, args ...interface{}){
	if logger.level <= ERROR_LEVEL {
		logger.nlogger.SetPrefix("Error ")
		logger.nlogger.Output(2, fmt.Sprintf(format, args...))
	}
}

func (logger *Logger) handlerFile() {
	ctime := time.Now()
	//设置的文件后缀，按照时间来拆分
	filenamePrefix := fmt.Sprintf("_%02d_%d_%2d",ctime.Year(),ctime.Month(),ctime.Day())

	os.MkdirAll(logger.fileFolder,os.ModePerm)

	clogFileName := fmt.Sprintf("%s%s",logger.fileName,filenamePrefix)

	// file exist, if not,should touch it
	_, err := isExist(clogFileName)
	if err!=nil {
		return
	}

	/**
	todo 支持大小分割
	 */
	fhandler, err := os.OpenFile(clogFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY,0777);
	if os.IsPermission(err) {
		fmt.Println(err)
		fmt.Println("permission deny")
		return
	}

	if  err==nil {
		logger.fileHandle.Sync()
		logger.fileHandle.Close()
		logger.fileHandle = fhandler
	}
}


func isExist(path string) (bool,error){
	_, err := os.Stat(path)
	if err==nil {
		return true,nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}