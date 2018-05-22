package pushcore

import (
	"gopkg.in/mgo.v2"
	"github.com/zuston/AtcalMq/util"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"time"
)

var session *mgo.Session

const (
	DATABASE  = "test"
	DB_PUSH_TAG = "queueTime"
)

func init(){
	session, err = mgo.Dial("127.0.0.1:27017")
	util.CheckPanic(err)
	session.SetMode(mgo.Monotonic, true)
}

func Mongo_Get(dbName string, convertMapper map[string]string, whereCondition string) []string {
	c := session.DB(DATABASE).C(dbName)

	var returnList []string
	var updateOidList []bson.ObjectId

	var tmp bson.M

	conditions := bson.M{
		"$or": []bson.M{
			bson.M{DB_PUSH_TAG: bson.M{"$gt": whereCondition}},
			bson.M{DB_PUSH_TAG: bson.M{"$exists":false}}}}

	iter := c.Find(conditions).Iter()
	//_ = conditions
	for iter.Next(&tmp){
		updateOidList = append(updateOidList,tmp["_id"].(bson.ObjectId))
		delete(tmp,"_id")
		// key convert
		tmpConvert := convert(tmp,convertMapper)
		jsonStr, _ := json.Marshal(tmpConvert)
		returnList = append(returnList,string(jsonStr))
	}

	updateTime := time.Now().Format("2006-01-02 15:04:05")
	// update time
	for _ ,oid := range updateOidList{
		data := bson.M{"$set":bson.M{DB_PUSH_TAG:updateTime}}
		err = c.UpdateId(oid,data)
		if !util.CheckErr(err){
			zlogger.Error("更新失败,oid:%s,err:%s",oid,err)
		}
	}

	return returnList
}
func convert(ms bson.M, convertMapper map[string]string) map[string]string {
	mapper := make(map[string]string,10)
	for key,value := range ms{
		correspondKey,ok := convertMapper[key]
		if ok {
			mapper[correspondKey] = value.(string)
		}else {
			mapper[key] = value.(string)
		}
	}
	return mapper
}
