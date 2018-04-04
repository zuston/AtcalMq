package util

import "strings"

/**
handle some json data
 */
func SplitByTag(jsonLine string, tag string) []string{
	arr := strings.Split(jsonLine,tag)
	return arr
}

