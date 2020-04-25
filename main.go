package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
)

func main(){

	c := colly.NewCollector(

		colly.AllowedDomains("www.qimao.com"),
	)
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

		fmt.Printf("名称：%s \n",tit)
		fmt.Printf("主角：%s \n",regstr[0][1])
		fmt.Printf("作者：%s \n",pname)
		fmt.Printf("状态：%s \n",status)
		fmt.Printf("数量：%s %s %s \n",nums,nums1,nums2)
		fmt.Printf("跟新时间：%s \n",uptime)
		fmt.Printf("最新章节：%s \n",newpage)
	})
	//获取url
	c.OnHTML(`a[href]`,func(e *colly.HTMLElement){
		link := e.Attr("href")
		reg := regexp.MustCompile(`^(/shuku/[1-9]{6}/)`)
		res := reg.FindAllString(link,-1)
		if len(res)>0{
			fmt.Printf("link Found:https://www.qimao.com%s \n",link)
		}


	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit("https://www.qimao.com/shuku/151822/")
}
