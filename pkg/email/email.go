// Package email implements the EmailSvc abstracting the platform we will be using to end our emails
package email

import (
	"context"
	"errors"

	"github.com/mailgun/mailgun-go/v4"
)

var (
	ErrEmptyEmail   = errors.New("empty email not allowed")
	ErrNoRecipients = errors.New("property 'To' in email was empty. No recipients provided")
)

type EmailSvc interface {
	Send(ctx context.Context, email *Email) (err error)
}

type MailgunOps struct {
	Domain string `json:"domain"`
	ApiKey string `json:"apiKey"`
}

type Email struct {
	TemplateId string
	Subject    string
	Html       string
	From       string
	Variables  *map[string]any
	To         []string
	Cc         []string
}

type MailgunSvc struct {
	client mailgun.MailgunImpl
}

// NewEmailService creates a new email service pointer
func NewMailgunSvc(ops *MailgunOps) EmailSvc {
	return &MailgunSvc{
		client: *mailgun.NewMailgun(ops.Domain, ops.ApiKey),
	}
}

// NewEmail
func NewEmail(templateId, html, subject, from string, variables *map[string]any, to, cc []string) *Email {
	return &Email{
		TemplateId: templateId,
		Html:       html,
		Subject:    subject,
		From:       from,
		Variables:  variables,
		To:         to,
		Cc:         cc,
	}
}

// Send triggers an email to be send using the provided email configuration.
func (s *MailgunSvc) Send(ctx context.Context, email *Email) (err error) {
	if email == nil || (email.Html == "" && email.TemplateId == "") {
		return ErrEmptyEmail
	}
	if email.To == nil || len(email.To) == 0 {
		return ErrNoRecipients
	}
	m := s.client.NewMessage(email.From, email.Subject, "", email.To...)

	if email.Html == "" {
		m.SetTemplate(email.TemplateId)

		for k, v := range *email.Variables {
			err = m.AddTemplateVariable(k, v)

			if err != nil {
				return err
			}
		}
	} else {
		m.SetHtml(email.Html)
	}
	_, _, err = s.client.Send(ctx, m)
	return err
}
