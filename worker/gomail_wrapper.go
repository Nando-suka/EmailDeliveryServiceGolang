package worker

import "gopkg.in/gomail.v2"

type GomailDialer interface {
	DialAndSend(m *gomail.Message) error
}

type RealDialer struct {
	Host string
	Port int
	User string
	Pass string
}

func (d *RealDialer) DialAndSend(m *gomail.Message) error {
	return gomail.NewDialer(d.Host, d.Port, d.User, d.Pass).DialAndSend(m)
}
