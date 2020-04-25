package tools

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type Gorm struct {
	*gorm.DB
}
//定义全局DB
var DB *Gorm

//数据连接返回全局DB
func InitGorm() (*Gorm,error){
	//读取配置
	database := Cfg.DatabaseConfig
	//定义连接字符
	args :=fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		database.User, database.Password,database.Host,database.Port,database.Database,database.Charset)
	//连接数据库
	gdb,err := gorm.Open(database.Driver,args)
	//设置全局变量赋值
	DB = &Gorm{gdb}
	if err !=nil{
		return nil,err
	}
	//返回全局变量
	return DB,nil
}
