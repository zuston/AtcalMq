package main

import (
	_ "github.com/zuston/AtcalMq/core"
	"fmt"
	"github.com/tidwall/gjson"
)

func main() {
	jsonline := `[{"name":"zuston","age":null,"period":1.23,"details":[{"old":46}]},{"name":"shacha"}]`
	//v := core.ModelGen([]byte(jsonline))
	//for _, value := range v{
	//	for key,va := range value{
	//		fmt.Println(key,string(va))
	//	}
	//}

	fmt.Println(gjson.Get(jsonline,"#.name"))

	result := gjson.Get(jsonline, ".#")
	result.ForEach(func(key, value gjson.Result) bool {
		fmt.Println(gjson.Get(value.String(),"name"))
		return true // keep iterating
	})
}