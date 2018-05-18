package main

import (
	"github.com/zuston/AtcalMq/util"
	"github.com/zuston/AtcalMq/rabbitmq"
	"time"
	"fmt"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"encoding/json"
	"io/ioutil"
	"github.com/zuston/AtcalMq/core/object"
	"strings"
	"flag"
)


var mysqlConn mysql.Conn
var logger *util.Logger

const(
	mysqlUsername = "root"
	mysqlPassword = "zuston"
)

const(
	MODEL_SUFFIX = "model"
)

const PUSH_LOG_PATH  = "/tmp/AnePush.log"

var(
	duration = flag.Duration("duration",time.Minute,"please enter the push duration")
	optionPath = flag.String("option","/opt/optional.model","the push queue column name")
	verifyTime = flag.String("verifyTime","","the push time")
)

func init(){
	flag.Parse()
	mysqlConn = mysql.New("tcp", "", "127.0.0.1:3306", mysqlUsername, mysqlPassword, "todo")
	mysqlConn.Register("set names utf8")

	err := mysqlConn.Connect()
	if err != nil {
		panic(err)
	}

	logger,err = util.NewLogger(util.DEBUG_LEVEL,PUSH_LOG_PATH)
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

	//duration := 10*time.Second
	//optionPath := "./optional.model"
	//verifyTime := "2018-09-09 23:00:00"
	timelyPublish(pf,*duration,*optionPath,*verifyTime)
}


func timelyPublish(pf *rabbitmq.ProducerFactory, duration time.Duration, optionPath string, verifyTime string) {
	translationT := ParseTranslationTable(optionPath)

	ticker := time.NewTicker(duration)

	for _ = range ticker.C{
		for componentKey,translationMapper := range translationT{
			// get the data from db and push the queue
			get2push(pf,componentKey,translationMapper,verifyTime)
		}
	}
}

func get2push(pf *rabbitmq.ProducerFactory, componentKey string, translationMapper map[string]string, verifyTime string) {
	tablename, queuename := parseComponentKey(componentKey)
	// 获取表的各column name
	columns := getMetaColumn(tablename)
	// 校验对应关系
	if !checkTranslation(translationMapper,columns) {
		logger.Error("[%s],配置文件与表的field不相互对应")
		panic("错误对应")
	}

	// 获取信息组合成 json
	pushInfos := getData2Json(tablename,columns,translationMapper,verifyTime)

	// 依次发送，或者分批次发送
	for _, pushInfo := range pushInfos{
		// 对应的tablename to queuename
		//pf.Publish(queuename,pushInfo)
		_ = queuename
		fmt.Println(pushInfo)
	}
}

func parseComponentKey(componentKey string) (string, string) {
	arr := strings.Split(componentKey,"_")
	return arr[0],arr[1]
}

func getData2Json(tablename string, columnsname []string, translationMapper map[string]string,verifyTime string) []string {
	// 增加对应的时间选择
	whereSql := "where queueTime is null "
	if verifyTime!="" {
		whereSql += fmt.Sprintf(`or queueTime>"%s"`,verifyTime)
	}
	// 从效率考虑，当然可以分段分批次取出，懒得写了
	// todo 有心人写一下呗
	conditionSql := fmt.Sprintf("select * from %s %s",tablename,whereSql)
	rows, _,err := mysqlConn.Query(conditionSql)
	checkErr(err)
	var retList []string
	// 需要更新推送时间的 主键 list
	var updateIdList []int
	for _,row := range rows{
		fieldsLen := len(row)
		// 全部设置为string
		rowMap := make(map[string]string)
		updateIdList = append(updateIdList,row.Int(0))
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

	idKeyName := columnsname[0]
	updateDataOfTimestamp(tablename,updateIdList,idKeyName)
	return retList
}

// 更新数据时间
func updateDataOfTimestamp(tablename string, ids []int,idKeyName string) {
	idWheres := zmap(ids, func(v int) string{
		return fmt.Sprintf("%s=%d",idKeyName,v)
	})
	whereSql := "where "+strings.Join(idWheres," OR ")
	updateTime := time.Now().Format("2006-01-02 15:04:05")
	updateSql := fmt.Sprintf("update %s set queueTime='%s' %s",tablename,updateTime,whereSql)

	mysqlConn.Query(updateSql)
}

// like the map reduce function
func zmap(ints []int, mapFunc func(v int) string) []string {
	var list []string
	for _,value := range ints{
		list = append(list,mapFunc(value))
	}
	return list
}


func checkTranslation(translationMapper map[string]string, columns []string) bool {
	for _,column := range columns{
		if _,ok := translationMapper[column]; !ok{
			fmt.Println(column)
			return false
		}
	}
	return true
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
func ParseTranslationTable(modelPath string) map[string]map[string]string {
	// 兼容传递目录或者直接是文件
	files := []string{modelPath}
	if util.IsDir(modelPath) {
		files = util.WalkDir(modelPath,MODEL_SUFFIX)
	}
	modelMappers := make(map[string]map[string]string)

	for _,file := range files{
		data, err := ioutil.ReadFile(file)
		checkPanic(err)
		datajson := []byte(data)

		var models []object.PushMapperObj
		err = json.Unmarshal(datajson, &models)
		checkPanic(err)


		for _,model := range models{
			tablename := model.TABLENAME
			queuename := model.QUEUENAME

			componentKey := fmt.Sprintf("%s_%s",tablename,queuename)
			modelMapper := make(map[string]string)
			for _,relation := range model.REALTIONS{
				modelMapper[relation.CN] = relation.QN
			}

			modelMappers[componentKey] = modelMapper
		}
	}
	return modelMappers
}

func checkPanic(err error) {
	if err!=nil {
		logger.Error("读取model映射文件报错")
		panic(err)
	}	
}

