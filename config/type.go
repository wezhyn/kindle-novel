package config

type Novel struct {
	Name   string `yaml:"name"`
	Author string `yaml:"author"`
	Cycle  bool   `yaml:"cycle"`
	Rules  []Rule `yaml:"rule"`
}

type Rule struct {
	Url string `yaml:"url"`
	Use string `yaml:"use"`
}
type Email struct {
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Port        int    `yaml:"port"`
	Host        string `yaml:"host"`
	Receiver    string `yaml:"receiver"`
	ErrReceiver string `yaml:"errReceiver"`
}

type Site struct {
	WebName  string `yaml:"name"`
	Method   string `yaml:"method"`
	Update   string `yaml:"update"`
	Refer    string `yaml:"refer"`
	Body     string `yaml:"body"`
	Next     string `yaml:"next"`
	NextName string `yaml:"nextName"`
}
