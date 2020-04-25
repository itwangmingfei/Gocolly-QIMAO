package tools

import (
	"github.com/monnand/goredis"
	"log"
)

type GoRedis struct {
	client goredis.Client
}
var Redis *GoRedis

const (
	DO_QUEUE = "DoQueue"//记录需要执行的数据
	TO_QUEUE = "Toqueue"//处理过的数据放这里
)
//链接redis
func InitRedis() *GoRedis{
	var cliect  goredis.Client
	cliect.Addr = Cfg.RedisConfig.Addr+":"+Cfg.RedisConfig.Port
	cliect.Password = Cfg.RedisConfig.Password
	cliect.Db = Cfg.RedisConfig.Db
	Redis = &GoRedis{client:cliect}
	return Redis
}
//存入数据头部存入数据
func (red *GoRedis) DoLpush(value string) error {
	err := red.client.Lpush(DO_QUEUE,[]byte(value))
	if err!=nil{
		log.Println(err)
		return err
	}
	return  nil
}
//读取数据尾部获取数据
func(red *GoRedis) DoRpop()(string,error){
	val,err := red.client.Rpop(DO_QUEUE)
	if err!=nil{
		return "",err
	}
	return string(val),err
}
//获取当前key存入的数量
func(red *GoRedis) DoLen() int{
	nums,err :=red.client.Llen(DO_QUEUE)
	if err!=nil {
		return 0
	}
	return nums
}

//存入集合中------------------------------------
func(red *GoRedis) ToSadd(value string) bool {
	res,_ :=red.client.Sadd(TO_QUEUE,[]byte(value))
	return res
}
//获取当前value是否需要存
func(red *GoRedis) ToIsset(value string) bool{
	res,_ :=red.client.Sismember(TO_QUEUE,[]byte(value))
	return res
}
//---------------------------------------------------

