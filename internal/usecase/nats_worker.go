package usecase

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/ImamTry257/Notify-Service/internal/entity"
	"github.com/ImamTry257/Notify-Service/internal/repository"
	"github.com/ImamTry257/Notify-Service/pkg/email"
)

//go:embed templates/*.html
var emailTemplates embed.FS

type NATSWorker struct {
	js        nats.JetStreamContext
	repo      repository.NotifyRepository
	email     email.SMTPSender
	templates *template.Template
	subject   string
	durable   string
}

func NewNATSWorker(js nats.JetStreamContext, repo repository.NotifyRepository, email email.SMTPSender, subject, durable string) *NATSWorker {
	// Pre-parse templates
	tmpl := template.Must(template.ParseFS(emailTemplates, "templates/*.html"))

	return &NATSWorker{
		js:        js,
		repo:      repo,
		email:     email,
		templates: tmpl,
		subject:   subject,
		durable:   durable,
	}
}

func (w *NATSWorker) Start(ctx context.Context) error {
	// ... (same as before)
	sub, err := w.js.PullSubscribe(w.subject, w.durable)
	if err != nil {
		return err
	}

	log.Printf("Worker started, listening on subject: %s", w.subject)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msgs, err := sub.Fetch(1, nats.Context(ctx))
				if err != nil {
					if err == context.DeadlineExceeded || err == nats.ErrTimeout {
						continue
					}
					log.Printf("Fetch error: %v", err)
					continue
				}

				for _, msg := range msgs {
					w.processMessage(ctx, msg)
				}
			}
		}
	}()

	return nil
}

func (w *NATSWorker) processMessage(ctx context.Context, msg *nats.Msg) {
	var history entity.EmailHistory
	if err := json.Unmarshal(msg.Data, &history); err != nil {
		log.Printf("Unmarshal error: %v", err)
		msg.Term()
		return
	}

	log.Printf("Processing notification for ID: %d, Email: %s, Type: %s", history.ID, history.Email, history.Type)

	// Send HTML email
	err := w.sendRealEmail(&history)
	
	status := "SENT"
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		status = "FAILED"
	}

	// Update record in MySQL
	if err := w.repo.UpdateStatus(ctx, history.ID, status); err != nil {
		log.Printf("Failed to update status in DB: %v", err)
	}

	msg.Ack()
}

func (w *NATSWorker) sendRealEmail(h *entity.EmailHistory) error {
	var templateName string
	var subject string

	switch h.Type {
	case entity.TYPE_OTP:
		templateName = "otp.html"
		subject = "Your Security OTP"
	case entity.TYPE_ACTIVATION:
		templateName = "activation.html"
		subject = "Activate Your Account"
	case entity.TYPE_RESET_PASSWORD:
		templateName = "reset_password.html"
		subject = "Password Reset Request"
	default:
		return fmt.Errorf("unsupported notification type: %s", h.Type)
	}

	var body bytes.Buffer
	if err := w.templates.ExecuteTemplate(&body, templateName, h); err != nil {
		return err
	}

	return w.email.SendHTMLEmail(h.Email, subject, body.String())
}

func (w *NATSWorker) sendEmail(h *entity.EmailHistory) error {
	// Deprecated simulated method
	return nil
}
