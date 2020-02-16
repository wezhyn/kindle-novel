package executor

import (
	"github.com/gocolly/colly"
	"strings"
)

/**
请求更新页面
*/
func (t *Tune) RequestUpdatePage(method string, url string, ctx *colly.Context) error {
	return t.Request(method, url, ctx, true)
}

/**
发出实际的请求
ctx: 为该请求特有的上下文内容
*/
func (t *Tune) CrawlerNovel(ctx *colly.Context) error {
	item := GetUpdateInfo(ctx)
	return t.Request("GET", item.Url, ctx, false)

}
func (t *Tune) CrawlerNovelNext(url string, ctx *colly.Context) error {
	return t.Request("GET", url, ctx, false)
}

/**
isUpdate: c.CrawlerCollector.Request
*/
func (t *Tune) Request(method string, url string, ctx *colly.Context, isUpdate bool) error {
	var m string
	if method == "" {
		m = "GET"
	} else {
		m = strings.ToUpper(method)
	}
	var c *colly.Collector
	// 该网站全局id
	// 为每次请求生成一个id，避免mian线程过早退出
	InitNovelId(ctx)
	t.wg.Add(1)
	if isUpdate {
		cache := GetRule(ctx)
		t.Info("%d :获取 %s 更新页面 %s \n", GetNovelId(ctx), cache.BookName, url)
		c = t.updateCollector
	} else {
		item := GetUpdateInfo(ctx)
		t.Info("%d: 获取 %s %s 第 %d 页 ： %s", GetNovelId(ctx), item.BookName, item.Title, GetPageNum(ctx), item.Url)
		c = t.crawlerCollector
	}
	return c.Request(m, url, nil, ctx, nil)
}
