package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	"gocoll/model"
	"gocoll/tools"
	"log"
	"regexp"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main(){
	//加载配置
	RedConfig()
	//初始化redis
	tools.InitRedis()
	//初始化Gorm
	tools.InitGorm()

	tools.DB.AutoMigrate(&model.Novel{})
	//*********************************
	c := colly.NewCollector(
		colly.AllowedDomains("www.qimao.com"),

		//这次在colly.NewCollector里面加了一项colly.Async(true)，表示抓取时异步的
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	///这个地方是配置redis存储请求返回的数据信息吧？
	storage := &redisstorage.Storage{
		Address: tools.Cfg.RedisConfig.Addr+":"+tools.Cfg.RedisConfig.Port,
		Password:tools.Cfg.RedisConfig.Password,
		DB: tools.Cfg.RedisConfig.Db,
		Prefix : "HTTP_QIMAO",
	}
	if err := c.SetStorage(storage);err!=nil{
		panic(err)
	}
	if err:=storage.Clear();err!=nil{
		log.Fatal(err)
	}
	//defer storage.Client.Close()
	//*----------------------------------------------
	var novel model.Novel
	q,_ :=queue.New(2,storage)
	/*<h2 class="tit">????</h2>*/
	c.OnHTML(`div.data-txt`, func(e *colly.HTMLElement) {
		ls,_ := e.DOM.Html()

		//标题
		tit := e.ChildText("h2.tit")
		//作者
		pname := e.ChildText(`p.p-name a`)
		//状态
		status := e.ChildText(`span.qm-tags.black.clearfix em:first-child`)

		nums := e.ChildText(`p.p-num span:nth-child(1)`)
		nums1 := e.ChildText(`p.p-num span:nth-child(3)`)
		nums2 := e.ChildText(`p.p-num span:nth-child(5)`)

		//没有标签正则一下
		reg := regexp.MustCompile(`<em>主角：</em>(.*?)<`)
		regstr := reg.FindAllStringSubmatch(ls,-1)

		uptime := e.ChildText(`p.p-update em.time`)

		newpage :=e.ChildText(`p.p-update a`)


		if len(tit)>0 {
			novel.Title = tit
			fmt.Printf("名称：%s \n",tit)
		}
		if len(regstr) >0 {
			novel.Mster = regstr[0][1]
			fmt.Printf("主角：%s \n",regstr[0][1])
		}
		if len(pname) >0 {
			fmt.Printf("作者：%s \n",pname)
			novel.Author = pname
		}
		if len(status) >0 {
			fmt.Printf("状态：%s \n",status)
			novel.Status = status
		}
		if len(nums)>0{
			fmt.Printf("数量：%s %s %s \n",nums,nums1,nums2)
			novel.Show = nums+nums1+nums2
		}
		if len(uptime)>0{
			fmt.Printf("跟新时间：%s \n",uptime)
			novel.Uptime =uptime
		}
		if len(newpage)>0{
			fmt.Printf("最新章节：%s \n",newpage)
			novel.Newpage = newpage
		}

		tools.DB.Create(&novel)
	})
	//获取url
	c.OnHTML(`a[href]`,func(e *colly.HTMLElement){
		link := e.Attr("href")
		reg := regexp.MustCompile(`^(/shuku/[1-9]{6}/)`)
		res := reg.FindAllString(link,-1)
		//存在返回数据
		if len(res)>0{
			//存入redis中
			newLink :=  "https://www.qimao.com" +link
			//判断执行存储的redis中是否存在如果存在不存储当前url
			Isset := tools.Redis.ToIsset(newLink)

			//如果不存在继续抓取当前链接
			if !Isset{
				//连接存入redis中
				err := tools.Redis.DoLpush(newLink)
				if err!=nil{
					log.Println(err)
				}
			}

		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	//获取这个应该是配饰redisstorage使用 用于接收返回数据然后存储
	c.OnResponse(func(r *colly.Response) {
		log.Println(c.Cookies(r.Request.URL.String()))
	})

	//处理请求url
	for {
		dlink,_ := tools.Redis.DoRpop()
		if len(dlink)==0{
			dlink ="https://www.qimao.com"
		}
		novel.Url = dlink
		q.AddURL(dlink)
		q.Run(c)
		time.Sleep(time.Second*2)
	}

}

func RedConfig(){
	cpath := "config/config.json"
	tools.ReadConfig(cpath)

}
