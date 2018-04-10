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

var strring = `

[{"areaCode":"B3","ewbNo":"300151337277","ewbsListNo":"49724048218013001","hewbNo":"30015133727700030002","nextSiteName":"漳州分拨中心","palletNo":"048255201802010552037450letTime":"2018-02-01 05:53:09","palletType":"1","palletUserCode":"048255","platformCode":"001","prevSiteName":"松江分拨中心","servicesType":"标准快运","siteCode":"0592001",048,"volume":0.33,"weight":76.67},{"areaCode":"B3","ewbNo":"300150837423","ewbsListNo":"49724048218013001","hewbNo":"30015083742300030001","nextSiteName":"漳州分拨中心","pa":"048255201802010552037450","palletTime":"2018-02-01 05:53:18","palletType":"1","palletUserCode":"048255","platformCode":"001","prevSiteName":"松江分拨中心","servicesType"","siteCode":"0592001","siteId":4048,"volume":0.14,"weight":28.33}]

`
var objectName = "CenterPalletObj"


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