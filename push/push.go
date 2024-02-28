package push

import (
	"net"
	"net/smtp"
)

type Push interface {
	push() error
}

type SmtpP struct {
	MyAddress string
	to        []string
	pwd       string
	smtpAddr  string
	//ibbattcwquxjbhbe
	auth smtp.Auth
}

func NewSmtpP(myAddress string, to []string, pwd string, smtpAddr string) *SmtpP {
	host, _, _ := net.SplitHostPort(smtpAddr)
	s := &SmtpP{MyAddress: myAddress, to: to, pwd: pwd, smtpAddr: smtpAddr}
	s.auth = smtp.PlainAuth("", s.MyAddress, s.pwd, host)
	return s

}

func (s *SmtpP) Push(msg []byte) error {

	return smtp.SendMail(s.smtpAddr, s.auth, s.MyAddress, s.to, msg)

}
