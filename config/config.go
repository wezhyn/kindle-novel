package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type Config struct {
	UserAgent  string   `yaml:"userAgent"`
	BodyIgnore []string `yaml:"bodyIgnore"`
	Site       []Site   `yaml:"site"`
	Email      *Email   `yaml:"email"`
	Novels     []Novel  `yaml:"novels"`
}

func (c *Config) GetSiteRule(site string) (Site, error) {
	for _, s := range c.Site {
		if s.WebName == site {
			return s, nil
		}
	}
	return Site{}, fmt.Errorf("无匹配的规则")
}

func (c *Config) UpdateUserAgent(agent string) {
	c.UserAgent = agent
}
func (c *Config) UpdateNovels(novels []Novel) {
	c.Novels = make([]Novel, len(novels), len(novels))
	copy(c.Novels, novels)
}

func (c *Config) GetLinkNum() int {
	sum := 0
	for _, n := range c.Novels {
		sum += len(n.Rules)
	}
	return sum
}

func (n *Novel) UpdateName(name string) {
	n.Name = name
}

func (n *Novel) UpdateAuthor(author string) {
	n.Author = author
}

func (n *Novel) UpdateLinks(links []Rule) {
	/*	n.Rules = make(map[string]string)
		for r, v := range links {
			n.Rules[r] = v
		}
	*/
	copy(n.Rules, links)
}

/**
用于定义config.yml 的地址的配置类
*/
type Path struct {
	homeDir    string
	configFile string
}

func (p *Path) ModifyHomeDir(homeDir string) {
	p.homeDir = abs(homeDir)
}
func (p *Path) ModifyConfigFile(configFile string) {
	p.configFile = configFile
}
func (p Path) ConfigFile() string {
	return path.Join(p.homeDir, p.configFile)
}

func abs(path string) string {
	if !filepath.IsAbs(path) {
		currentDir, _ := os.Getwd()
		return filepath.Join(currentDir, path)
	}
	return path
}
