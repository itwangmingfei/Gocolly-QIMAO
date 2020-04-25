package model

type Novel struct {
	Title string `gorm:"commit:'标题'"`
	Mster string  `gorm:"commit:'主角'"`
	Author string  `gorm:"commit:'作者'"`
	Status string `gorm:"commit:'状态'"`
	Show string  `gorm:"commit:'数量统计'"`
	Uptime string `gorm:"commit:'更新时间'"`
	Newpage string  `gorm:"commit:'最新章节'"`
	Url string `gorm:"commit:'URL'"`
}
