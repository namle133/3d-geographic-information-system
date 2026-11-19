package services

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/uit/vinhomes-management/config"
)

type PasswordResetSender interface {
	Configured() bool
	Send(recipient, token string) error
}

type SMTPPasswordResetSender struct{ cfg *config.Config }

func NewSMTPPasswordResetSender(cfg *config.Config) *SMTPPasswordResetSender {
	return &SMTPPasswordResetSender{cfg: cfg}
}
func (s *SMTPPasswordResetSender) Configured() bool {
	return s.cfg != nil && s.cfg.SMTPHost != "" && s.cfg.SMTPUser != "" && s.cfg.SMTPPassword != "" && s.cfg.SMTPFrom != ""
}
func (s *SMTPPasswordResetSender) Send(recipient, token string) error {
	if !s.Configured() {
		return errors.New("password reset email is not configured")
	}
	if strings.ContainsAny(recipient+s.cfg.SMTPFrom, "\r\n") {
		return errors.New("invalid email address")
	}
	to, err := mail.ParseAddress(recipient)
	if err != nil {
		return err
	}
	from, err := mail.ParseAddress(s.cfg.SMTPFrom)
	if err != nil {
		return err
	}
	address := net.JoinHostPort(s.cfg.SMTPHost, s.cfg.SMTPPort)
	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}
	tlsConfig := &tls.Config{ServerName: s.cfg.SMTPHost, MinVersion: tls.VersionTLS12}
	if s.cfg.SMTPPort == "465" {
		conn = tls.Client(conn, tlsConfig)
	}
	client, err := smtp.NewClient(conn, s.cfg.SMTPHost)
	if err != nil {
		return err
	}
	defer client.Close()
	if s.cfg.SMTPPort != "465" {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return errors.New("SMTP server must support TLS")
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return err
		}
	}
	if err := client.Auth(smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)); err != nil {
		return err
	}
	if err := client.Mail(from.Address); err != nil {
		return err
	}
	if err := client.Rcpt(to.Address); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Dat lai mat khau Vinhomes\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nMã đặt lại mật khẩu của bạn:\r\n%s\r\n\r\nMã có hiệu lực 30 phút và chỉ dùng một lần. Nếu bạn không yêu cầu, hãy bỏ qua email này.\r\n", from.String(), to.String(), token)
	if _, err := writer.Write([]byte(message)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	// DATA has been accepted; a disconnect during QUIT does not undo delivery.
	_ = client.Quit()
	return nil
}
