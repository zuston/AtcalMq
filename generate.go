package main

import (
	"strings"
	"github.com/zuston/AtcalMq/util"
	"fmt"
	"log"
)

/**
generate the code
 */

var strring = `{"dataType":2,"ewbNo":"300179844074","ewbsListNo":"332114393218033103","hewbNo":"30017984407400040003","nextSiteCode":"5122002","nextSiteId":14393,"operatorCode":"002977,043056,00113317","platformCode":"006","scanMan":"002977","scanTime":"2018-04-01 02:03:39","scanType":1,"siteCode":"5103032","siteId":3321,"volume":0.18,"weight":35},`


var objectName = "CenterLoadObj"


//var strring string
//var objectName string



func main(){

	//fmt.Println(fmt.Sprintf("%c[1;40;32m[Please input your json]%c[0m",0x1B,0x1B))
	//fmt.Scanln(&strring)
	//fmt.Println(fmt.Sprintf("%c[1;40;32mPlease input your want the structName%c[0m",0x1B,0x1B))
	//fmt.Scanln(&objectName)
	snippetJson := getSnippet(strring)

	componentLine := strings.Replace(snippetJson,"{","",1)
	componentLine = strings.Replace(componentLine,"}","",1)

	log.Println(componentLine)

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
		if strings.Contains(value,`"`) {
			objectValue = "string"
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


	fmt.Println(fmt.Sprintf("%c[1;40;32mGENERATE HBASE MAPPER AS FOLLOWING%c[0m",0x1B,0x1B))

	// generate the mapper hbase code
	mapperFirstLine := "map[string][]byte{"
	mapperEndLine := "}"



	fmt.Println(mapperFirstLine)
	for k,v := range objectMapper{
		if v=="float32" {
			fmt.Println(fmt.Sprintf(`"%s":[]byte(fmt.sprintf("\\\%s",v.%s)),`,k,k))
		}
		fmt.Println(fmt.Sprintf(`"%s":[]byte(v.%s),`,k,k))
	}

	fmt.Println(mapperEndLine)

 }


 func getSnippet(jsons string) string{
	firstIndex := strings.Index(jsons,"{")

	secondIndex := strings.Index(jsons[firstIndex:],"}")

	return jsons[firstIndex:secondIndex+firstIndex+1]
 }