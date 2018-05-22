package pushcore

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"fmt"
	"encoding/json"
	"strings"
	"time"
	"github.com/zuston/AtcalMq/util"
)


var(
	mysqlConn mysql.Conn
)

const(
	mysqlDB = "todo"
	mysqlUsername = "root"
	mysqlPassword = "zuston"
)


func init(){
	mysqlConn = mysql.New("tcp", "", "127.0.0.1:3306", mysqlUsername, mysqlPassword, mysqlDB)
	mysqlConn.Register("set names utf8")
	err := mysqlConn.Connect()
	util.CheckPanic(err)
}

func Mysql_Get(dbName string, convertMapper map[string]string, whereCondition string) []string {
	// 获取表的各column name
	columns := getMetaColumn(dbName)
	// 校验对应关系
	if !checkConvertRelation(convertMapper,columns) {
		zlogger.Error("[%s],配置文件与表的field不相互对应")
		panic("错误对应")
	}
	return mysql_getData2Json(dbName,columns,convertMapper,whereCondition)
}

// 获取 mysql table 元信息
func getMetaColumn(tablename string) []string {
	var columnNameList []string
	metaColumnSql := fmt.Sprintf("DESCRIBE %s",tablename)
	rows,_,_ := mysqlConn.Query(metaColumnSql)
	for _,row := range rows{
		columnNameList = append(columnNameList,row.Str(0))
	}
	return columnNameList
}

// 校验获取到的信息与转换表的关系
func checkConvertRelation(translationMapper map[string]string, columns []string) bool {
	for _,column := range columns{
		if _,ok := translationMapper[column]; !ok{
			fmt.Println(column)
			return false
		}
	}
	return true
}

func mysql_getData2Json(tablename string, columnsname []string, translationMapper map[string]string,verifyTime string) []string {
	// 增加对应的时间选择
	whereSql := "where queueTime is null "
	if verifyTime!="" {
		whereSql += fmt.Sprintf(`or queueTime>"%s"`,verifyTime)
	}
	// 从效率考虑，当然可以分段分批次取出，懒得写了
	// todo 有心人写一下呗
	conditionSql := fmt.Sprintf("select * from %s %s",tablename,whereSql)
	rows, _,err := mysqlConn.Query(conditionSql)
	util.CheckErr(err)
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
	// 结果为0,无需更新
	if len(retList)==0 {
		return retList
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
