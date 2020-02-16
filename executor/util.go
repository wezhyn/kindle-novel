package executor

import (
	"hash/fnv"
	"math/rand"
	"time"
)

func characterId(g *GatherData) (uniqueId int) {
	uniqueId = Hashcode(g.NovelRule.BookName) + g.Number
	return uniqueId

}

func CreateIdentifyId(name string) int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	nameCode := int64(Hashcode(name))
	return r.Int63() + nameCode
}
func Hashcode(str string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(str))
	return int(h.Sum32())
}
