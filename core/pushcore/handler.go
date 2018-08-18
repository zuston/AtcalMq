package pushcore

import (
	"github.com/zuston/AtcalMq/core"
	"log"
	"github.com/zuston/AtcalMq/util"
)

const(
	LOGGER_PATH = "/extdata/log/AnePushHandler.log"
)

var(
	zlogger *util.Logger
	err error
)

func init(){
	zlogger,err = util.NewLogger(util.DEBUG_LEVEL,LOGGER_PATH)
	util.CheckPanic(err)
	zlogger.SetDebug()
}

/**
处理push逻辑
 */
func Get(dbContainerTag int,dbName string,convertMapper map[string]string,whereCondition string) []string{

	switch dbContainerTag {
	case core.MONGO_TAG:
		return Mongo_Get(dbName,convertMapper,whereCondition)
	case core.MYSQL_TAG:
		return Mysql_Get(dbName,convertMapper,whereCondition)
	default:
		log.Println("default none")
	}
	return nil
}



