package executor

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"kindle/data"
	M "kindle/data/mysql"
	P "kindle/parse"
	"reflect"
	"strings"
	"time"
)

type GatherData struct {
	IdentifyId int64
	// 当前章节的url
	DataUrl string
	// 当前小说章节
	Number     int
	Next       bool
	Title      string
	PageNum    int
	Body       []byte
	CreateTime int64
	NovelRule  *NovelRule
}

/**
获取实际的小说章节
*/
func HandleNovelChapter(response *colly.Response, t Tune) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(response.Body))
	rule := GetRule(response.Ctx)
	item := GetUpdateInfo(response.Ctx)
	if err != nil {
		panic(err)
	}
	body := HandleBody(doc, rule)
	nextTag := doc.Find(rule.Next)
	pageNum := GetPageNum(response.Ctx)
	var isNext bool
	// 获取下一页
	if nextTag != nil {
		nextName := nextTag.Text()
		if nextName == rule.NextName {
			attr, exists := nextTag.Attr("href")
			if exists {
				url := response.Request.AbsoluteURL(attr)
				go func() {
					context := CopyContext(response.Ctx)
					HandleContextAddPage(context)
					HandleRetryResponseError(3, context, &t, t.CrawlerNovelNext, url, context)
				}()
				isNext = true
			}
		}
	}
	gatherData := GatherData{
		Next:       isNext,
		Title:      item.Title,
		PageNum:    pageNum,
		DataUrl:    item.Url,
		Number:     item.Number,
		Body:       body,
		NovelRule:  rule,
		CreateTime: time.Now().Unix(),
		IdentifyId: GetNovelId(response.Ctx),
	}
	t.SendChannelData(&gatherData)
}

// 获取当前页内容
func HandleBody(doc *goquery.Document, rule *NovelRule) []byte {
	htmlBody := doc.Find(rule.Body)
	var body bytes.Buffer
	body.Write([]byte("<br>"))
	for _, con := range htmlBody.Contents().Nodes {
		text := con.Data
		text = strings.TrimPrefix(text, "\n")
		text = strings.TrimSpace(text)
		if text != "" && len(text) >= 2 {
			var ign = false
			for _, ignore := range rule.BodyIgnore {
				if ignore == text {
					ign = true
				}
			}
			if !ign {
				text = strings.TrimSpace(text)
				text = "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;" + text + "<br>"
				body.Write([]byte(text))
				body.Write([]byte("<br>"))
			}
		}
	}
	return body.Bytes()
}

/**
  处理更新页面，例如：https://m.xbiquge.cc/chapters_rev_50627/1
  返回 UpdateItem 列表
*/
func HandleUpdatePage(response colly.Response) data.UpdateItems {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(response.Body))
	rule := GetRule(response.Ctx)
	if err != nil {
		panic(err)
	}
	items := make(data.UpdateItems, 0, 16)
	refers := strings.Split(rule.Refer, ";")
	referNode := refers[0]
	referAttribute := refers[1]
	doc.Find(rule.Update).EachWithBreak(func(i int, selection *goquery.Selection) bool {
		if i > 10 {
			return false
		}
		attr, exists := selection.Find(referNode).Attr(referAttribute)
		headline := selection.Text()
		item, err := P.Title(headline)
		if exists && err == nil {
			item.Url = response.Request.AbsoluteURL(attr)
			item.UrlHash = Hashcode(item.Url)
			item.BookName = rule.BookName
			item.WebName = rule.WebName
			item.Cycle = rule.Cycle
			items = append(items, *item)
		}
		return true
	})
	return items
}

/**
处理更新列表，判断哪些是未更新的章节
*/
func HandleUpdateItems(items []data.UpdateItem, response colly.Response) []data.UpdateItem {
	bookService := M.NewInstance()
	rule := GetRule(response.Ctx)
	cycle := rule.Cycle
	currentNum, _ := bookService.Last(rule.BookName)
	strategy := HandleItemStrategy(currentNum, cycle)
	return strategy.Select(items)
}

/**
获取策略
	* currentNum==-1 代表第一次更新，获取第一个
 	* currentNum!=-1 cycle=false ，获取大于 currentNum 的Item
	* currentNum!=-1 cycle=true , 2，1，47，46，45 时
		* currentNum= 46 获取 2，1，47
		* currentNum= 1, 获取 2
*/
func HandleItemStrategy(num int, cycle bool) data.ItemStrategy {
	if num == -1 {
		return data.NoFetchStrategy{}
	} else {
		if cycle {
			return data.CycleFetchStrategy{CurrentNum: num}
		} else {
			return data.MaxFetchStrategy{CurrentNum: num}
		}
	}
}
func HandleRetryResponseError(retry int, ctx *colly.Context, t *Tune, f interface{}, args ...interface{}) {
	v := reflect.ValueOf(f)
	if v.Kind() != reflect.Func {
		panic("错误的传递类型，请检查方法调用")
	}
	ft := reflect.TypeOf(f)
	argv := make([]reflect.Value, ft.NumIn())
	for i := 0; i < ft.NumIn(); i++ {
		if xt, targ := reflect.TypeOf(args[i]), ft.In(i); xt != nil && !xt.AssignableTo(targ) {
			panic(fmt.Errorf("reflect: f using  %s as Type %s", xt.String(), targ.String()))
		} else if xt != nil {
			argv[i] = reflect.ValueOf(args[i])
		}
	}
	if retry <= 0 {
		retry = 3
	}
	for i := 0; i < retry; i++ {
		result := v.Call(argv)[0].Interface()
		if result != nil {
			err := result.(error)
			if taskId := GetNovelIdNative(ctx); taskId != nil && err != nil {
				t.UnRegister(taskId.(int64))
				t.Error("error: %s", err.Error())
			}
		} else {
			break
		}
	}
}
