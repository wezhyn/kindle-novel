package parse

import (
	"fmt"
	"github.com/Wall-ee/chinese2digits/Go/chinese2digits"
	"hash/fnv"
	"kindle/data"
	"regexp"
	"strconv"
)

func Title(tag string) (*data.UpdateItem, error) {
	handleTitle := chinese2digits.TakeChineseNumberFromString(tag).(map[string]interface{})
	numberTitle := handleTitle["replacedText"]
	re := regexp.MustCompile(`第(\d{1,4})章.(.*)`)
	subMatchs := re.FindStringSubmatch(numberTitle.(string))
	it := data.UpdateItem{}
	var err error = fmt.Errorf("标题错误")
	if len(subMatchs) == 3 {
		it.Title = subMatchs[2]
		it.Number, err = strconv.Atoi(subMatchs[1])
	}
	return &it, err
}

func hashcode(str string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(str))
	return int(h.Sum32())
}
