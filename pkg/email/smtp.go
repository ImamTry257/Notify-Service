package email

import (
	"fmt"
	"net/smtp"

	"github.com/ImamTry257/Notify-Service/config"
)

type SMTPSender interface {
	SendHTMLEmail(to string, subject string, body string) error
}

type smtpSender struct {
	cfg config.SMTPConfig
}

func NewSMTPSender(cfg config.SMTPConfig) SMTPSender {
	return &smtpSender{cfg: cfg}
}

func (s *smtpSender) SendHTMLEmail(to string, subject string, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	
	// Message headers
	header := make(map[string]string)
	header["From"] = fmt.Sprintf("%s <%s>", s.cfg.SenderName, s.cfg.SenderEmail)
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Auth (only if username is provided)
	var auth smtp.Auth
	if s.cfg.Username != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}

	return smtp.SendMail(addr, auth, s.cfg.SenderEmail, []string{to}, []byte(message))
}
