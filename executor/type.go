package executor

import "kindle/config"

/**
爬取网站的结构，在response与request之间共享
*/

type NovelRule struct {
	// 书本名
	BookName string
	// 网站名
	WebName string

	// 更新章节的网站
	UpdateUrl string
	// 作者
	Author string
	// 循环型小说
	Cycle bool

	// 当前小说更新进度
	CurrentNumber int
	// 读取Body时忽略的字段
	BodyIgnore []string
	config.Site
}
