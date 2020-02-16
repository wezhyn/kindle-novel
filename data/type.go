package data

import (
	"bytes"
	"hash/fnv"
)

/**
用于 kindle/parse 用于解析标题
*/
type UpdateItem struct {
	BookName string
	WebName  string
	Title    string
	//章节
	Number  int
	Url     string
	Cycle   bool
	UrlHash int
}
type UpdateItems []UpdateItem

func (u UpdateItems) Len() int {
	return len(u)
}

func (u UpdateItems) Less(i, j int) bool {
	priorityI := getCyclePriority(u[i].Number)
	priorityJ := getCyclePriority(u[j].Number)
	return priorityI-priorityJ > 0
}

func (u UpdateItems) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type TxtItem struct {
	UpdateItem
	Body *bytes.Buffer
}

type TxtItems []TxtItem

func (t TxtItems) Len() int {
	return len(t)
}

func (t TxtItems) Less(i, j int) bool {
	ti := t[i]
	tj := t[j]
	var ni, nj int = ti.Number, tj.Number
	if tj.Cycle {
		nj = getCyclePriority(tj.Number)
	}
	if ti.Cycle {
		ni = getCyclePriority(ti.Number)
	}
	vi := ni + hashcode(ti.BookName)
	vj := nj + hashcode(tj.BookName)
	// 倒序
	return vi-vj < 0
}

func (u TxtItems) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func hashcode(str string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(str))
	return int(h.Sum32())
}

type IBookData interface {
	number() int
	name() string
	title() string
	url() string
	send() bool
}
