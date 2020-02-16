package executor

import (
	"bytes"
	"fmt"
	"github.com/gocolly/colly"
	"kindle/config"
	"kindle/data"
	E "kindle/data/email"
	"kindle/data/mysql"
	"kindle/data/txt"
	log2 "log"
	"os"
	"sort"
	"sync"
)

type Tune struct {
	// 用于检查小说是否更新
	updateCollector *colly.Collector
	//	爬取更新的小说
	crawlerCollector *colly.Collector
	channel          chan *GatherData
	wg               *sync.WaitGroup
	// 收集每一个爬虫任务的 IdentifyId,用来错误时删除任务
	// true: 未删除
	works       *sync.Map
	mutex       *sync.Mutex
	infoLogger  *log2.Logger
	errorLogger *log2.Logger
}

func newCollector(c config.Config) *colly.Collector {
	return colly.NewCollector(
		colly.UserAgent(c.UserAgent),
		colly.DetectCharset(),
		colly.AllowURLRevisit(),
	)
}

func NewInstance(c config.Config) *Tune {
	t := new(Tune)
	t.updateCollector = newCollector(c)
	t.updateCollector.Async = true
	t.crawlerCollector = newCollector(c)
	t.channel = make(chan *GatherData, 3)
	t.works = new(sync.Map)
	t.mutex = &sync.Mutex{}
	t.wg = new(sync.WaitGroup)
	t.infoLogger = log2.New(os.Stdout, "INFO: ", log2.Ldate|log2.Ltime|log2.Lshortfile)
	t.errorLogger = log2.New(os.Stdout, "ERROR: ", log2.Ldate|log2.Ltime|log2.Lshortfile)
	go t.receiveChannelData()
	t.updateCollector.OnResponse(func(response *colly.Response) {
		t.Info("%d -1 获取更新链接成功 ", GetNovelId(response.Ctx))
		t.UnRegister(GetNovelId(response.Ctx))
		pageItems := HandleUpdatePage(*response)
		sort.Sort(pageItems)
		updateItems := HandleUpdateItems(pageItems, *response)
		if updateItems != nil && len(updateItems) > 0 {
			for _, item := range updateItems {
				context := CopyContext(response.Ctx)
				// 初始化 1
				PutPageNum(context, 1)
				PutUpdateInfo(context, &item)
				HandleRetryResponseError(3, context, t, t.CrawlerNovel, context)
			}
		} /*else {
			t.UnRegister(getNovelId(response.Ctx))
		}*/
	})
	t.crawlerCollector.OnResponse(func(response *colly.Response) {
		HandleNovelChapter(response, *t)
	})
	return t
}

func (t Tune) Info(format string, a ...interface{}) {
	t.infoLogger.Printf(fmt.Sprintf(format, a...))
}

func (t *Tune) Error(format string, a ...interface{}) {
	t.errorLogger.Printf(fmt.Sprintf(format, a...))
}

func (t *Tune) UnRegister(task int64) {
	t.mutex.Lock()
	store, loaded := t.works.LoadOrStore(task, false)
	if loaded {
		// 加载已读
		v := store.(bool)
		if v {
			t.works.Store(task, false)
			t.Info("taskId : %d -1", task)
			t.wg.Done()
		}
	} else {
		//	 存取
		t.Info("taskId : %d -1", task)
		t.wg.Done()
	}
	t.mutex.Unlock()
}

/**
只存一次，其他时候只读
*/
func (t *Tune) Register(task int64) bool {
	store, loaded := t.works.LoadOrStore(task, true)
	if loaded {
		//	已存在，返回是否取消
		return store.(bool)
	} else {
		return true
	}
}

/**
用于收集小说页面
*/
var buffer sync.Map
var alreadBuffer sync.Map

/**
单线程运作
*/
func (t *Tune) receiveChannelData() {
	novel := t.channel
	for g := range novel {
		// 计算当前小说 title+ Number 唯一标识符
		uniqueId := characterId(g)
		gatherData, exist := buffer.Load(uniqueId)
		// 判断当前章节是否已经读取
		_, loaded := alreadBuffer.Load(uniqueId)
		if !loaded && !exist {
			// 未读取并未有值
			gatherData = make([]GatherData, 0, 8)
		}
		if !loaded {
			gatherData = append(gatherData.([]GatherData), *g)
			t.handleNovelComplete(gatherData.([]GatherData))
		}
	}
}

/**
检查每章节是否完成,完成后添加到已读取列表
未完成
*/
func (t *Tune) handleNovelComplete(novels []GatherData) {
	var isComplete bool = false
	for _, rgd := range novels {
		if rgd.Next == false && len(novels) == rgd.PageNum {
			// 已完成
			isComplete = true
			_, loaded := alreadBuffer.LoadOrStore(characterId(&rgd), true)
			if !loaded {
				// 存取当前 gataerData 为已读取
				t.togetherData(novels)
				buffer.Delete(characterId(&rgd))
			}
			//	 忽略已被读取的章节
		}
	}
	if !isComplete {
		// 未完成
		buffer.Store(characterId(&novels[0]), novels)
	}
}

var novelBuffer []data.TxtItem = make([]data.TxtItem, 0, 4)

/**
具体章节的处理
非阻塞
*/
func (t *Tune) togetherData(datas []GatherData) {
	resultData := new(data.TxtItem)
	resultData.Body = new(bytes.Buffer)
	for i, d := range datas {
		if i == 0 {
			resultData.WebName = d.NovelRule.WebName
			resultData.BookName = d.NovelRule.BookName
			resultData.Url = d.DataUrl
			resultData.Title = d.Title
			resultData.Cycle = d.NovelRule.Cycle
			resultData.Number = d.Number
			resultData.UrlHash = Hashcode(d.DataUrl)
		}
		resultData.Body.Write(d.Body)
	}
	t.Info("收集到 %s %s:%d", resultData.BookName, resultData.Title, resultData.Number)
	novelBuffer = append(novelBuffer, *resultData)
}

func (t *Tune) Wait() {
	t.wg.Wait()
	if len(novelBuffer) > 0 {
		clean := txt.Clean{}
		mobiPath, err := clean.Clear(novelBuffer)
		if err != nil {
			panic(err)
		}
		m := mysql.NewInstance()
		E.Send(mobiPath)
		os.Remove(mobiPath)
		for _, d := range novelBuffer {
			if err := m.Save(d.UpdateItem); err != nil {
				E.SendMsg(err.Error())
			}
		}
	}
}
func (t Tune) SendChannelData(novel *GatherData) {
	if t.Register(novel.IdentifyId) {
		t.Info("task :%d -1", novel.IdentifyId)
		t.wg.Done()
		t.channel <- novel
	}
}
