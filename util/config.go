package util

import (
	"os"
	"bufio"
	"strings"
	"io"
	"gopkg.in/ini.v1"
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
		if err!= nil {
			if err==io.EOF {
				return mapper, nil
			}
			return nil, err
		}
		line = strings.TrimSpace(line)
		tagIndex := strings.Index(line,"=")
		if tagIndex==-1 {
			return nil,err
		}
		key := strings.TrimSpace(line[:tagIndex])
		value := strings.TrimSpace(line[tagIndex+1:])
		mapper[key] = value

	}
}

// 读取 ini 配置文件，config path && sectionName
func NewConfigReader(path string,section string)(map[string]string, error){
	cfg, err := ini.Load(path)
	CheckPanic(err)
	_, err = cfg.GetSection(section)
	if err!=nil{
		return nil,err
	}
	containerMapper := make(map[string]string)
	allKeys := cfg.Section(section).KeyStrings()
	for _,key := range allKeys{
		containerMapper[key] = cfg.Section(section).Key(key).String()
	}
	return containerMapper,nil
}