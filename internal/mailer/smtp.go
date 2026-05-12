package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"kibit/internal/config"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	cfg *config.Config
	ch  chan *gomail.Message
}

func New(cfg *config.Config) *Mailer {
	slog.Info("📧 Initializing mailer",
		"host", cfg.SMTPHost,
		"port", cfg.SMTPPort,
		"user", cfg.SMTPUser,
		"from", cfg.SMTPFrom)

	m := &Mailer{
		cfg: cfg,
		ch:  make(chan *gomail.Message, 100),
	}
	go m.worker()
	return m
}

func (m *Mailer) SendAsync(ctx context.Context, to, subject, body string) {
	slog.Info("📨 Queuing email", "to", to, "subject", subject)

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.cfg.SMTPFrom)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	select {
	case m.ch <- msg:
		slog.Info("✅ Email added to queue", "to", to)
	default:
		slog.Error("❌ Mailer queue full, message dropped", "to", to)
	}
}

func (m *Mailer) worker() {
	slog.Info("🔄 Mailer worker started")

	port, _ := strconv.Atoi(m.cfg.SMTPPort)
	d := gomail.NewDialer(m.cfg.SMTPHost, port, m.cfg.SMTPUser, m.cfg.SMTPPass)
	slog.Debug("SMTP dialer created", "host", m.cfg.SMTPHost, "port", 587)

	for msg := range m.ch {
		slog.Info("🚀 Processing email from queue")
		if err := m.retrySend(d, msg); err != nil {
			slog.Error("❌ Email send failed after retries", "error", err)
		} else {
			slog.Info("✅ Email sent successfully")
		}
	}
}

func (m *Mailer) retrySend(d *gomail.Dialer, msg *gomail.Message) error {
	for i := 1; i <= 3; i++ {
		slog.Info("🌐 SMTP connect attempt", "attempt", i, "max", 3)

		err := d.DialAndSend(msg)
		if err == nil {
			slog.Info("✅ SMTP send successful", "attempt", i)
			return nil
		}

		slog.Warn("⚠️ SMTP error, will retry",
			"attempt", i,
			"error", err,
			"next_retry_in", time.Duration(i)*time.Second)

		time.Sleep(time.Duration(i) * time.Second)
	}

	return fmt.Errorf("failed after 3 attempts: last error logged above")
}
