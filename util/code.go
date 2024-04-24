package util

import (
	"context"
	"net/smtp"
	"strings"

	"github.com/chenjie199234/account/config"

	"github.com/chenjie199234/Corelib/util/common"
)

func SendEmailCode(ctx context.Context, email, code, action string) error {
	c := config.AC.Service
	host := ""
	if index := strings.LastIndex(c.EmailSMTPAddr, ":"); index == -1 {
		host = c.EmailSMTPAddr
	} else {
		host = c.EmailSMTPAddr[:index]
	}
	msg := []string{"From: " + c.EmailAccount, "To: " + email, "Subject: Login Dynamic Code\r\n", "Your dynamic code is:" + code}
	msgstr := strings.Join(msg, "\r\n")
	if host == c.EmailSMTPAddr {
		//use the default port 587
		return smtp.SendMail(c.EmailSMTPAddr+":587", smtp.PlainAuth("", c.EmailAccount, c.EmailPassword, host), c.EmailAccount, []string{email}, common.STB(msgstr))
	}
	return smtp.SendMail(c.EmailSMTPAddr, smtp.PlainAuth("", c.EmailAccount, c.EmailPassword, host), c.EmailAccount, []string{email}, common.STB(msgstr))
}

// TODO
func SendTelCode(ctx context.Context, tel, code, action string) error {
	return nil
}
