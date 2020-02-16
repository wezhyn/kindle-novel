package executor

import (
	"github.com/gocolly/colly"
	"kindle/data"
)

/**
用于处理 Context 的一些操作
*/

func PutRule(ctx *colly.Context, rule *NovelRule) {
	ctx.Put("cache", rule)
}
func CheckRule(ctx *colly.Context) bool {
	cc := ctx.GetAny("cache")
	return cc != nil
}
func GetRule(ctx *colly.Context) *NovelRule {
	return ctx.GetAny("cache").(*NovelRule)
}

func PutUpdateInfo(ctx *colly.Context, item *data.UpdateItem) {
	ctx.Put("item", item)
}
func GetUpdateInfo(ctx *colly.Context) *data.UpdateItem {
	return ctx.GetAny("item").(*data.UpdateItem)
}

/**
存取当前爬取到第几页
*/
func PutPageNum(ctx *colly.Context, pageNum int) {
	ctx.Put("page", pageNum)
}
func HandleContextAddPage(context *colly.Context) {
	pageNum := GetPageNum(context) + 1
	context.Put("page", pageNum)
}
func GetPageNum(context *colly.Context) int {
	pageNum := context.GetAny("page")
	if pageNum == nil {
		pageNum = 1
	}
	return pageNum.(int)
}

/**
每一个更新链接生成一个 Id，用于取消爬取任务
*/
func InitNovelId(ctx *colly.Context) {
	if !CheckRule(ctx) {
		panic("再初始化小说id前，先调用 putRule () 方法")
	}
	rule := GetRule(ctx)
	id := CreateIdentifyId(rule.UpdateUrl)
	ctx.Put("id", id)
}
func GetNovelId(ctx *colly.Context) int64 {
	return ctx.GetAny("id").(int64)
}
func GetNovelIdNative(ctx *colly.Context) interface{} {
	return ctx.GetAny("id")
}

func CopyContext(src *colly.Context) (desc *colly.Context) {
	desc = colly.NewContext()
	src.ForEach(func(k string, v interface{}) interface{} {
		desc.Put(k, v)
		return v
	})
	return desc
}
