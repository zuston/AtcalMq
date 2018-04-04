package util

import (
	"os"
	"bufio"
	"strings"
	"io"
)

// read the config file , convert to the Map[string]string
func ConfigReader(path string) (map[string]string, error){
	mapper := make(map[string]string,4)
	ffile, err := os.Open(path)
	if err!=nil {
		return nil, err
	}
	buffer := bufio.NewReader(ffile)
	for {
		line, err := buffer.ReadString('\n')
		line = strings.TrimSpace(line)
		tagIndex := strings.Index(line,"=")
		key := strings.TrimSpace(line[:tagIndex])
		value := strings.TrimSpace(line[tagIndex+1:])
		mapper[key] = value
		if err!= nil {
			if err==io.EOF {
				return mapper, nil
			}
			return nil, err
		}
	}
}
