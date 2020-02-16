package email

import (
	"gopkg.in/gomail.v2"
)

var (
	username    string = ""
	password    string = ""
	port        int    = 0
	host        string = ""
	receiver    string = ""
	errReceiver string = ""
)

func PostInit(un string, pw string, pt int, hs string, res string, errRes string) {
	username = un
	password = pw
	port = pt
	host = hs
	receiver = res
	errReceiver = errRes
}
func SendMsg(msg string) {
	Check()
	m := newMsgInstance()
	m.SetHeader("To", errReceiver)
	m.SetBody("text/html", msg)
	d := newDialer()
	//邮件发送服务器信息,使用授权码而非密码
	_ = d.DialAndSend(m)
}

func Send(attachPath string) {
	Check()
	m := newMsgInstance()
	m.Attach(attachPath)
	m.SetBody("text/html", "小说更新")
	d := newDialer()
	//邮件发送服务器信息,使用授权码而非密码
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
func newDialer() *gomail.Dialer {
	d := gomail.NewDialer(host, port, username, password)
	return d
}

func newMsgInstance() *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", username)  //发件人
	m.SetHeader("To", receiver)    //收件人
	m.SetHeader("Subject", "小说更新") //邮件标题
	return m
}
func Check() {
	if username == "" || password == "" || port == 0 || receiver == "" || host == "" {
		panic("请初始化邮箱配置")
	}
}
