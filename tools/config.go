package tools

import (
	"bufio"
	"encoding/json"
	"os"
)

type (
	Config struct {
		RedisConfig RedisConfi `json:"redis_config"`
		DatabaseConfig DatabaseConfig `json:"database_config"`

	}
	DatabaseConfig struct {
		Driver string `json:"driver"`
		User string `json:"user"`
		Password string `json:"password"`
		Host string `json:"host"`
		Port string `json:"port"`
		Database string `json:"database"`
		Charset string `json:"charset"`
		Sqlstr bool `json:"sqlstr"`
	}
	RedisConfi struct {
		Addr string `json:"addr"`
		Port string `json:"port"`
		Password string `json:"password"`
		Db int `json:"db"`
	}
)
var Cfg *Config = nil
//读取文件
func ReadConfig(path string)(*Config,error)  {
	file,err := os.Open(path)
	if err!=nil{
		panic(err)
	}
	//关闭文件
	defer file.Close()
	//读取文件
	readers := bufio.NewReader(file)
	//转义json
	decode :=json.NewDecoder(readers)
	//赋值
	if err = decode.Decode(&Cfg);err!=nil{
		return nil,err
	}
	return Cfg,nil
}