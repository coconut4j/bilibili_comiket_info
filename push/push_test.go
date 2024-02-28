package push

import (
	"fmt"
	"net/smtp"
	"testing"
)

func TestName(t *testing.T) {
	//s := NewSmtpP("caihualin96@aliyun.com", []string{"441483724@qq.com", "505174715@qq.com"}, "caihualin1996", "smtp.aliyun.com:25")

	// 创建邮件消息对象

	smtp_server := "smtp.aliyun.com"
	smtp_port := "25"

	senders_email := "caihualin96@aliyun.com"
	senders_password := "caihualin1996"

	recipient_email := "441483724@qq.com"
	message := []byte("To: " + recipient_email + "\r\n" +
		"Subject: Go SMTP Test\r\n" +
		"\r\n" +
		"Hello,\r\n" +
		"This is a test email sentfrom Go!\r\n"//接收退出信号退出程序mtp.PlainAuth("", senders_email, senders_password, smtp_server)

	err := smtp.SendMail(smtp_server+":"+smtp_port, auth, senders_email, []string{recipient_email}, message)
	if err != nil {
		panic(err)
	}
	fmt.Println("Email sent successfully!")

}
