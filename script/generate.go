package script

import (
	"strings"
	"github.com/zuston/AtcalMq/util"
	"fmt"
	"github.com/zuston/AtcalMq/core"
	"time"
)

/**
generate the code
 */

var strring = `
{"modifiedTime":"2018-04-10 12:33:45","viaCount":3,"lineName":"南京-合肥-六安1","driveName":"谢庭久"}
`
var objectName = "VehicleLineObj"


//var strring string
//var objectName string



func main(){

	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))

	for _,v := range core.BasicInfoTableNames{
		fmt.Println(fmt.Sprintf(`count "%s"`,v))
	}
	return

	//fmt.Println(fmt.Sprintf("%c[1;40;32m[Please input your json]%c[0m",0x1B,0x1B))
	//fmt.Scanln(&strring)
	//fmt.Println(fmt.Sprintf("%c[1;40;32mPlease input your want the structName%c[0m",0x1B,0x1B))
	//fmt.Scanln(&objectName)

	generateMultiObj(strring,objectName)
 }


 func generateMultiObj(jsonLine string,objectName string){
	 sonFirstIndex := strings.Index(strring,"[{")
	 if sonFirstIndex==-1 {
		 generateSimpleObj(jsonLine,objectName)
		 return
	 }
	 sonEndIndex := strings.Index(strring,"}]")
	 fmt.Println(sonFirstIndex,sonEndIndex)
	 fmt.Println(strring[sonFirstIndex:sonEndIndex+2])

	 sonObjName := ""
	 prefixLine := strring[:sonFirstIndex]
	 tempIndex := strings.LastIndex(prefixLine,",")
	 sonObjName = prefixLine[tempIndex+2:len(prefixLine)-2]+"Obj"

	 oneLine := strring[:sonFirstIndex]+fmt.Sprintf(`"[]%s"`,sonObjName)+strring[sonEndIndex+1:]

	 generateSimpleObj(oneLine,objectName)
	 generateSimpleObj(strring[sonFirstIndex:sonEndIndex+1],sonObjName)
 }


func generateSimpleObj(jsonLine string,objectName string) {
	snippetJson := getSnippet(jsonLine)

	componentLine := strings.Replace(snippetJson,"{","",1)
	componentLine = strings.Replace(componentLine,"}","",1)

	list := util.SplitByTag(componentLine,`,"`)

	objectMapper := make(map[string]string)
	for k,v := range list{
		if k==0 {
			v = v[1:]
		}
		temp := util.SplitByTag(v,":")
		key := temp[0]
		value := temp[1]

		objectKey := strings.ToUpper(key[:1]) + key[1:len(key)-1]
		objectValue := ""
		if strings.Contains(value,`"`) || value=="null" {
			if strings.Contains(value,"[]") {
				objectValue = value[1:len(value)-2]
			}else {
				objectValue = "string"
			}
		}else{
			if strings.Contains(value,".") {
				objectValue = "float32"
			}else{
				objectValue = "int"
			}
		}

		objectMapper[objectKey] = objectValue
	}

	firstLine := fmt.Sprintf("type %s struct {",objectName)

	fmt.Println(fmt.Sprintf("%c[1;40;32mGENERATE CODE AS FOLLOWING%c[0m",0x1B,0x1B))
	fmt.Println(firstLine)

	for k,v := range objectMapper{
		fmt.Println(k,v)
	}

	fmt.Println("}")


	//fmt.Println(fmt.Sprintf("%c[1;40;32mGENERATE HBASE MAPPER AS FOLLOWING%c[0m",0x1B,0x1B))
	//
	//// generate the mapper hbase code
	//mapperFirstLine := "map[string][]byte{"
	//mapperEndLine := "}"
	//
	//
	//
	//fmt.Println(mapperFirstLine)
	//for k,v := range objectMapper{
	//	if v=="float32" {
	//		fmt.Println(fmt.Sprintf(`"%s":[]byte(fmt.sprintf("\\\%s",v.%s)),`,k,k))
	//	}
	//	fmt.Println(fmt.Sprintf(`"%s":[]byte(v.%s),`,k,k))
	//}
	//
	//fmt.Println(mapperEndLine)
}


 func getSnippet(jsons string) string{
	firstIndex := strings.Index(jsons,"{")

	secondIndex := strings.LastIndex(jsons[firstIndex:],"}")

	return jsons[firstIndex:secondIndex+firstIndex+1]
 }


//count "ane_its_ai_data_centerLoad_queue"  						  16131493 row(s) in 909.1600 seconds
//count "ane_its_ai_biz_order_queue"                                752757 row(s) in 39.7860 seconds
//count "ane_its_ai_biz_ewb_queue"                                  1647591 row(s) in 98.3720 seconds
//count "ane_its_ai_basic_route_queue"                              102394 row(s) in 5.9590 seconds
//count "ane_its_ai_biz_ewbsList_queue"                             861840 row(s) in 43.1440 seconds
//count "ane_its_ai_data_siteLoad_queue"                            1753053 row(s) in 82.2560 seconds
//count "ane_its_ai_data_centerUnload_queue"                        16381973 row(s) in 816.6330 seconds
//count "ane_its_ai_data_centerPallet_queue"                        16064140 row(s) in 943.8790 seconds
//count "ane_its_ai_data_centerSort_queue"                          2184217 row(s) in 100.7350 seconds
//count "ane_its_ai_basic_area_queue"                               0 row(s) in 0.0490 seconds
//count "ane_its_ai_basic_attend_queue"                             0 row(s) in 0.0300 seconds
//count "ane_its_ai_trigger_siteSend_queue"                         0 row(s) in 0.0290 seconds
//count "ane_its_ai_trigger_siteUpload_queue"                       0 row(s) in 0.0280 seconds
//count "ane_its_ai_trigger_inOrOut_queue"                          0 row(s) in 0.0300 seconds
//count "ane_its_ai_trigger_stayOrLeave_queue"                      0 row(s) in 0.0290 seconds
//count "ane_its_ai_trigger_centerTransport_queue"                  657050 row(s) in 40.4940 seconds
//count "ane_its_ai_basic_site_queue"                               3715 row(s) in 0.3110 seconds
//count "ane_its_ai_basic_vehicleLine_queue"                        115222 row(s) in 6.5870 seconds
//count "ane_its_ai_basic_platform_queue"                           0 row(s) in 0.0310 seconds
