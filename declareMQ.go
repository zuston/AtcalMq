package main

import (
	"github.com/zuston/AtcalMq/core"
	"fmt"
)

func main(){
	fmt.Println("./rabbitmqadmin declare exchange --vhost=its-pds name=ane_its_exchange type=direct")
	for _, queueName := range core.BasicInfoTableNames{
		fmt.Printf("./rabbitmqadmin declare queue --vhost=its-pds name=%s durable=true",queueName)
		fmt.Println()
		fmt.Printf(`./rabbitmqadmin --vhost="its-pds" declare binding source="ane_its_exchange" destination_type="queue" destination="%s" routing_key="%s"`,queueName,queueName)
		fmt.Println()

		queueName = fmt.Sprintf("%s_%s",core.SAVING_RELATION_PREFIX,queueName)
		fmt.Printf("./rabbitmqadmin declare queue --vhost=its-pds name=%s durable=true",queueName)
		fmt.Println()

		fmt.Printf(`./rabbitmqadmin --vhost="its-pds" declare binding source="ane_its_exchange" destination_type="queue" destination="%s" routing_key="%s"`,queueName,queueName)
		fmt.Println()
		fmt.Println("#===========================================")

	}
}