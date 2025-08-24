package notifier


import (
"context"
"crypto/tls"
"fmt"
"net/smtp"
"strings"


"scansched/internal/config"
)


type email struct { cfg config.EmailConfig }


func NewEmail(c config.EmailConfig) Notifier { return &email{cfg: c} }


func (e *email) Notify(ctx context.Context, ev Event) error {
	if len(e.cfg.To) == 0 { return nil }
		auth := smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.SMTPHost)
		msg := fmt.Sprintf("Subject: %s\r\nFrom: %s\r\nTo: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s\r\n",
		ev.Title, e.cfg.From, strings.Join(e.cfg.To, ","), ev.Text)
		addr := fmt.Sprintf("%s:%d", e.cfg.SMTPHost, e.cfg.SMTPPort)
	if e.cfg.StartTLS {
		tlsconfig := &tls.Config{ServerName: e.cfg.SMTPHost}
		conn, err := tls.Dial("tcp", addr, tlsconfig)
	if err != nil { return err }
		c, err := smtp.NewClient(conn, e.cfg.SMTPHost)
	if err != nil { return err }
		defer c.Close()
	if err := c.Auth(auth); err != nil { return err }
	if err := c.Mail(e.cfg.From); err != nil { return err }
	for _, r := range e.cfg.To { if err := c.Rcpt(r); err != nil { return err } }
		w, err := c.Data(); if err != nil { return err }
	if _, err := w.Write([]byte(msg)); err != nil { return err }
	return w.Close()
}
	return smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, []byte(msg))
}