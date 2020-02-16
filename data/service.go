package data

/**
book 服务
*/
type BookService interface {
	// 获取书籍的最近一次更新章节，-1代表未获取
	Last(bookName string) (int, error)
	BookSave
}

type BookSave interface {
	Save(data UpdateItem) error
}
type BookClean interface {
	Clear(datas []TxtItem) error
}

type ItemStrategy interface {
	// 挑选策略
	// items 为一组根据 item.Number 排序后的数据
	Select(items []UpdateItem) []UpdateItem
}

type NoFetchStrategy struct {
}
type MaxFetchStrategy struct {
	CurrentNum int
}
type CycleFetchStrategy struct {
	CurrentNum int
}

/**
数据例如：2，1，46，45，44
CurrentNum=45
获取 2，1，46
*/
func (c CycleFetchStrategy) Select(items []UpdateItem) []UpdateItem {
	for i, item := range items {
		if item.Number == c.CurrentNum {
			return items[0:i]
		}
	}
	return items
}

/**
数据例如：47，46，45，44
CurrentNum=45 ,获取47，46
*/
func (m MaxFetchStrategy) Select(items []UpdateItem) []UpdateItem {
	for i, item := range items {
		if item.Number <= m.CurrentNum {
			return items[0:i:i]
		}
	}
	return items
}

func (n NoFetchStrategy) Select(items []UpdateItem) []UpdateItem {
	return items[0:1]
}
