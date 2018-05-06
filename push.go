package main

import (
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/rabbitmq"
	"time"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	"encoding/json"
)


var mysqlConn mysql.Conn
var logger *util.Logger

var tablename2queue map[string]string

const(
	mysqlUsername = "root"
	mysqlPassword = "shacha"
)

func init(){
	mysqlConn = mysql.New("tcp", "", "127.0.0.1:3306", mysqlUsername, mysqlPassword, "todo")
	mysqlConn.Register("set names utf8")

	err := mysqlConn.Connect()
	if err != nil {
		panic(err)
	}

	logger,err = util.NewLogger(util.DEBUG_LEVEL,"/tmp/push.log")
	if err!=nil {
		panic("init the log error")
	}
	logger.SetDebug()
}

// 定时向对方推送计算好的数据
func main(){
	// 1. 初始化producer配置
	// 2. 读取抓取推送的配置文件
	// 3. 设定定时器，定时推送
	// 4. 推送信息

	configMapper, _ := util.ConfigReader("/opt/mq.cfg")
	mq_uri := configMapper["mq_uri"]
	exchange := configMapper["exchange"]
	exchange_type := configMapper["exchange_type"]
	pf, _ := rabbitmq.NewProducerFactory(mq_uri,exchange,exchange_type,false)
	go pf.Handle()

	duration := 10*time.Second
	optionPath := "./xml.cfg"
	verifyTime := "2018-09-09 23:00:00"
	timelyPublish(pf,duration,optionPath,verifyTime)
}


func timelyPublish(pf *rabbitmq.ProducerFactory, duration time.Duration, optionPath string, verifyTime string) {
	translationT := parseTraslationTable(optionPath)

	ticker := time.NewTicker(duration)

	for _ = range ticker.C{
		for tableName,translationMapper := range translationT{
			// get the data from db and push the queue
			get2push(pf,tableName,translationMapper,verifyTime)
		}
	}
}

func get2push(pf *rabbitmq.ProducerFactory, tablename string, translationMapper map[string]string, verifyTime string) {
	// 获取表的各column name
	columns := getMetaColumn(tablename)
	// 校验对应关系
	if !checkTranslation(translationMapper,columns) {
		logger.Error("[%s],配置文件与表的field不相互对应")
		panic("错误对应")
	}

	// 获取信息组合成 json
	pushInfos := getData2Json(tablename,columns,translationMapper)
	// 依次发送，或者分批次发送
	for _, pushInfo := range pushInfos{
		// 对应的tablename to queuename
		pf.Publish(tablename2queue[tablename],pushInfo)
	}
}

func getData2Json(tablename string, columnsname []string, translationMapper map[string]string) []string {
	conditionSql := fmt.Sprintf("select * from %s where ptime is null or ptime >= '%s'",tablename,verifyTime)
	rows, _,err := mysqlConn.Query(conditionSql)
	checkErr(err)
	var retList []string
	for _,row := range rows{
		fieldsLen := len(row)
		// 全部设置为string
		rowMap := make(map[string]string)
		for i:=0;i<fieldsLen;i++{
			columnName := columnsname[i]
			// 两者要完全一致才行
			correspondQueueName,ok := translationMapper[columnName]
			if ok {
				rowMap[correspondQueueName] = row.Str(i)
			}else {
				rowMap[columnName] = row.Str(i)
			}
		}
		convertJson, _ := json.Marshal(rowMap)
		retList = append(retList,string(convertJson))
	}
	return retList
}


func checkTranslation(strings map[string]string, strings2 []string) bool {
	return false
}

func getMetaColumn(tablename string) []string {
	var columnNameList []string
	metaColumnSql := fmt.Sprintf("DESCRIBE %s",tablename)
	rows,_,_ := mysqlConn.Query(metaColumnSql)
	for _,row := range rows{
		columnNameList = append(columnNameList,row.Str(0))
	}
	return columnNameList
}

func checkErr(e error) {
	if e!=nil {
		logger.Error("err : [%s]",e)
	}
}

// 形式为 ：{"tablename":{"name":"zname","age":"zage"}}
func parseTraslationTable(s string) map[string]map[string]string {
	return nil
}

// 定时推送

