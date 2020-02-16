package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"testing"
)

func Test_read(t *testing.T) {
	c := Config{}
	//c.Novels = make([]Novel, 4)
	if file, err := ioutil.ReadFile("../config.yml"); err != nil {
		panic(err)
	} else {
		_ = yaml.Unmarshal(file, &c)
	}
	fmt.Println(c)
}
