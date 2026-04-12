package usecase

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/ImamTry257/Notify-Service/internal/entity"
	"github.com/ImamTry257/Notify-Service/internal/repository"
	"github.com/ImamTry257/Notify-Service/pkg/email"
	"github.com/ImamTry257/Notify-Service/pkg/logger"
	"github.com/nats-io/nats.go"
)

const (
	flowNATSWorkerStart   = "NATSWorker.Start"
	flowNATSWorkerProcess = "NATSWorker.ProcessMessage"
	flowNATSWorkerEmail   = "NATSWorker.SendEmail"
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

	logger.Info(flowNATSWorkerStart, "worker started", "subject", w.subject)

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
					logger.Error(flowNATSWorkerStart, "fetching messages", err, "subject", w.subject)
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
		logger.Error(flowNATSWorkerProcess, "unmarshaling message", err)
		msg.Term()
		return
	}

	logger.Request(flowNATSWorkerProcess, history)

	// Send HTML email
	err := w.sendRealEmail(&history)

	status := "SENT"
	if err != nil {
		logger.Error(flowNATSWorkerProcess, "sending email", err,
			"id", history.ID,
			"email", history.Email,
			"type", history.Type,
		)
		status = "FAILED"
	} else {
		logger.Response(flowNATSWorkerProcess, map[string]any{
			"id": history.ID, "email": history.Email, "type": history.Type, "status": status,
		})
	}

	// Update record in MySQL
	if err := w.repo.UpdateStatus(ctx, history.ID, status); err != nil {
		logger.Error(flowNATSWorkerProcess, "updating status in DB", err,
			"id", history.ID,
			"status", status,
		)
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
