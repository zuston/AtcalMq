package pullcore

import (
	"github.com/zuston/AtcalMq/core"
	"fmt"
	"time"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"github.com/zuston/AtcalMq/util"
)

var relationMappers map[int][]string

const (
	RELATION_A = iota
	RELATION_B
	RELATION_C
)

const TIME_FORMAT = "2006-01-02 15:04:05"

func init(){
	relationMappers = map[int][]string{
		RELATION_A : []string{
			core.QUEUE_BIZ_EWB,
			core.QUEUE_BASIC_ROUTE,
		},
		RELATION_B : []string{
			core.QUEUE_DATA_SITELOAD,
			core.QUEUE_DATA_CENTERUNLOAD,
			core.QUEUE_DATA_CENTERPALLET,
			core.QUEUE_DATA_CENTERSORT,
			core.QUEUE_DATA_CENTERLOAD,
		},
		RELATION_C : []string{
			core.QUEUE_BIZ_EWBSLIST,
			core.QUEUE_DATA_SITELOAD,
			core.QUEUE_DATA_CENTERUNLOAD,
			core.QUEUE_DATA_CENTERPALLET,
			core.QUEUE_DATA_CENTERLOAD,
		},
	}
}

/**
 存储多重队列的关系，此处为二元关系组
 */
func MultiSavingModel(qn string, info map[string]string) func(){
	// 三张表的是否关联的判断
	table_A, table_B, table_C := relationMatch(qn)

	return func() {

		if table_A {
			// hbase table name
			table_A_NAME := ""

			// rowkey 由 startSiteID 和 ewbNo 组成
			startSiteId := info["startsiteid"]
			ewbNo := info["ewbno"]
			rowkey := fmt.Sprintf("%s#%s",startSiteId,ewbNo)
			// 扩充 qualifier
			info["mqDate"] = time.Now().Format(TIME_FORMAT)

			// 存储形式
			columnFamily := qn
			qualifierMapper := map[string]map[string][]byte{
				columnFamily : V2Byte(info),
			}

			// rpc save action
			saveAction(table_A_NAME,rowkey,qualifierMapper)
		}

		if table_B {
			table_B_NAME := ""
			siteId := info["siteid"]
			ewbNo := info["ewbno"]

			rowkey := fmt.Sprintf("%s#%s#%s",siteId,time.Now().Format(TIME_FORMAT),ewbNo)

			info["mqDate"] = time.Now().Format(TIME_FORMAT)
			cf := qn
			qmapper := map[string]map[string][]byte{
				cf : V2Byte(info),
			}

			saveAction(table_B_NAME,rowkey,qmapper)
		}

		if table_C {
			table_C_NAME := ""
			siteId := info["siteid"]
			ewbNo := info["ewbno"]

			rowkey := fmt.Sprintf("%s#%s#%s",siteId,time.Now().Format(TIME_FORMAT),ewbNo)

			info["mqDate"] = time.Now().Format(TIME_FORMAT)
			cf := qn
			qmapper := map[string]map[string][]byte{
				cf : V2Byte(info),
			}

			saveAction(table_C_NAME,rowkey,qmapper)
		}
	}
}

func saveAction(tableName string, rowkey string, qualifierMapper map[string]map[string][]byte) {
	putReq, _ := hrpc.NewPutStr(context.Background(),tableName,rowkey,qualifierMapper)
	_, err := singleConn.Client.Put(putReq)
	if !util.CheckErr(err) {

	}
}

/**
	拆分函数
 */
func relationMatch(queueName string) (bool, bool, bool) {
	a, b, c := false,false,false
	for key, values := range relationMappers{
		if key==RELATION_A && contains(values,queueName) {
			a = true
		}

		if key==RELATION_B && contains(values,queueName) {
			b = true
		}

		if key==RELATION_C && contains(values,queueName) {
			c = true
		}
	}
	return a, b, c
}

func contains(list []string, qn string) bool{
	for _,value := range list{
		if value==qn {
			return true
		}
	}
	return false
}
