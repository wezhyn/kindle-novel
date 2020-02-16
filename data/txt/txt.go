package txt

import (
	"github.com/766b/mobi"
	"io"
	"kindle/data"
	E "kindle/data/email"
	log2 "log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	coverMakingHost string = "https://orly.nanmu.me"
	coverMakingPath string = "/api/generate?"
)

func PostInit(coverMaking string) {
	coverMakingHost = coverMaking
}

type Clean struct {
}

/**
将获取的文本清洗成kindle接受的格式
*/
func (c Clean) Clear(datas data.TxtItems) (string, error) {
	mobiName := "./" + strconv.FormatInt(time.Now().Unix(), 16) + ".mobi"
	mobiWriter, err := mobi.NewWriter(mobiName)
	mobiPath, _ := filepath.Abs(mobiName)
	if err != nil {
		E.SendMsg("无法创建 mobi 文件")
		panic("Error: can't init mobi %s" + err.Error())
	}
	mobiWriter.Compression(mobi.CompressionNone) // LZ77 compression is also possible using  mobi.CompressionPalmDoc
	title := mobiTitle()
	mobiWriter.Title(title)
	coverpath := downCover()
	defer os.Remove(coverpath)
	mobiWriter.AddCover(coverpath, coverpath)
	// Meta data
	mobiWriter.NewExthRecord(mobi.EXTH_DOCTYPE, "PDOC")
	mobiWriter.NewExthRecord(mobi.EXTH_AUTHOR, "Kindle")
	mobiWriter.NewExthRecord(mobi.EXTH_PUBLISHER, "Wezhyn")
	sort.Sort(datas)
	for _, d := range datas {
		mobiWriter.NewChapter(d.Title, d.Body.Bytes())
	}
	mobiWriter.Write()
	log2.Print("Info: 创建 %s 成功 \n", strings.TrimSpace(title))
	return mobiPath, nil
}

func downCover() string {
	var coverUrl url.URL
	tm := time.Now()
	q := coverUrl.Query()
	guideTitle := tm.Format("Jan pm Monday ")
	imageId := tm.Unix() % 40
	title := mobiTitle()
	q.Add("g_loc", "BR")
	q.Add("g_text", guideTitle)
	q.Add("color", "70706d")
	q.Add("img_id", strconv.FormatInt(imageId, 10))
	q.Add("author", "Kindle")
	q.Add("top_text", "加速死亡")
	q.Add("title", title)

	for i := 0; i < 5; i++ {
		requestUrl := coverMakingHost + coverMakingPath + q.Encode()
		resp, err := http.Get(requestUrl)
		if err == nil && resp.StatusCode == 200 {
			file, err := os.Create("./" + strings.TrimSpace(title) + ".jpg")
			if err != nil {
				panic(err)
			}
			p, _ := filepath.Abs(file.Name())
			defer file.Close()
			_, copyErr := io.Copy(file, resp.Body)
			if copyErr != nil {
				panic(copyErr)
			}
			return p
		}
	}
	defaultUrl, _ := filepath.Abs("./default.jpg")
	return defaultUrl
}

func mobiTitle() (title string) {
	h := time.Now().Hour()
	htimes := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	index := (h + 1) / 2
	title = htimes[index] + "报"
	return "        " + title
}
