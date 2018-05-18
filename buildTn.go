package main

import (
	"github.com/zuston/AtcalMq/core"
	"fmt"
)

const (
	Max = 1 << 30
)

func main(){
	for _,v := range core.BasicInfoTableNames{
		createStr := fmt.Sprintf(`create "%s",{NAME=>"%s"}`,v,"basic")
		fmt.Println(createStr)
	}

	for _,tn := range core.LinkTableNames{
		headerStr := fmt.Sprintf(`create "%s", `,tn)
		arr := []string{}
		for _, cftn := range core.BasicInfoTableNames{
			createStr := fmt.Sprintf(`{NAME=>"%s"},`,cftn)
			arr = append(arr,createStr)
		}
		if tn=="Link_Site" {
			a := fmt.Sprintf(`{NAME=>"%s"},`,"nextSite_ane_its_ai_biz_ewbsList_queue")
			b := fmt.Sprintf(`{NAME=>"%s"},`,"nextSite_ane_its_ai_data_siteLoad_queue")
			c := fmt.Sprintf(`{NAME=>"%s"},`,"nextSite_ane_its_ai_data_centerLoad_queue")

			arr = append(arr,a,b,c)
		}
		fmt.Println(headerStr,fmt.Sprintf("%s",arr)[1:])
	}
}
