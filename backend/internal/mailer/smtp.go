package mailer

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type SMTP struct {
	Address  string
	Username string
	Password string
	From     string
	Timeout  time.Duration
}

func (sender SMTP) SendPasswordReset(ctx context.Context, recipient, resetURL string) error {
	from, err := mail.ParseAddress(sender.From)
	if err != nil {
		return fmt.Errorf("parse reset sender: %w", err)
	}
	to, err := mail.ParseAddress(recipient)
	if err != nil {
		return fmt.Errorf("parse reset recipient: %w", err)
	}
	host, _, err := net.SplitHostPort(sender.Address)
	if err != nil {
		return fmt.Errorf("parse SMTP address: %w", err)
	}
	timeout := sender.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	dialer := net.Dialer{Timeout: timeout}
	connection, err := dialer.DialContext(ctx, "tcp", sender.Address)
	if err != nil {
		return fmt.Errorf("connect SMTP: %w", err)
	}
	defer connection.Close()
	_ = connection.SetDeadline(time.Now().Add(timeout))

	client, err := smtp.NewClient(connection, host)
	if err != nil {
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer client.Close()
	if supported, _ := client.Extension("STARTTLS"); !supported {
		return errors.New("SMTP server does not support STARTTLS")
	}
	if err := client.StartTLS(&tls.Config{MinVersion: tls.VersionTLS12, ServerName: host}); err != nil {
		return fmt.Errorf("start SMTP TLS: %w", err)
	}
	if sender.Username != "" {
		if err := client.Auth(smtp.PlainAuth("", sender.Username, sender.Password, host)); err != nil {
			return fmt.Errorf("authenticate SMTP: %w", err)
		}
	}
	if err := client.Mail(from.Address); err != nil {
		return fmt.Errorf("set reset sender: %w", err)
	}
	if err := client.Rcpt(to.Address); err != nil {
		return fmt.Errorf("set reset recipient: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("open SMTP body: %w", err)
	}
	message := resetMessage(from.String(), to.String(), resetURL)
	if _, err := io.Copy(writer, strings.NewReader(message)); err != nil {
		return fmt.Errorf("write reset email: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close reset email: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("finish SMTP delivery: %w", err)
	}
	return nil
}

func resetMessage(from, to, resetURL string) string {
	var builder strings.Builder
	writer := bufio.NewWriter(&builder)
	_, _ = fmt.Fprintf(writer, "From: %s\r\n", from)
	_, _ = fmt.Fprintf(writer, "To: %s\r\n", to)
	_, _ = fmt.Fprint(writer, "Subject: Restablece tu acceso a APU\r\n")
	_, _ = fmt.Fprint(writer, "MIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n")
	_, _ = fmt.Fprint(writer, "Auto-Submitted: auto-generated\r\n\r\n")
	_, _ = fmt.Fprint(writer, "Recibimos una solicitud para restablecer tu contraseña de APU.\r\n\r\n")
	_, _ = fmt.Fprintf(writer, "Abre este enlace durante los próximos 30 minutos:\r\n%s\r\n\r\n", resetURL)
	_, _ = fmt.Fprint(writer, "Si no solicitaste el cambio, ignora este mensaje. El enlace solo puede usarse una vez.\r\n")
	_ = writer.Flush()
	return builder.String()
}
