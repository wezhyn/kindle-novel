package executor

import (
	"github.com/gocolly/colly"
	"kindle/config"
)

func Fun(c config.Config) error {
	t := NewInstance(c)
	// 每一个小说可能有多个更新地址
	for _, novel := range c.Novels {
		// 一个更新地址附加更新细节
		for _, link := range novel.Rules {
			site, ruleErr := c.GetSiteRule(link.Use)
			if ruleErr != nil {
				panic(ruleErr)
			}
			rule := NovelRule{
				BookName:      novel.Name,
				WebName:       link.Use,
				Author:        novel.Author,
				UpdateUrl:     link.Url,
				Cycle:         novel.Cycle,
				BodyIgnore:    c.BodyIgnore,
				CurrentNumber: -1,
				Site:          site,
			}
			ctx := colly.NewContext()
			PutRule(ctx, &rule)
			HandleRetryResponseError(3, ctx, t, t.RequestUpdatePage, rule.Method, rule.UpdateUrl, ctx)
		}
	}
	t.Wait()
	return nil
}
