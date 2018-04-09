package core

import (
	"strings"
	"fmt"
)

const (
	CONSUMER_TAG_PREFIX = "CONSUMER-"
)

func GetQueueNameFromConsumerTag(s string) string {
	index := strings.Index(s,"-")
	return s[index+1:]
}
func SetConsumerTag(queueName string) string{
	return fmt.Sprintf("%s%s",CONSUMER_TAG_PREFIX,queueName)
}
