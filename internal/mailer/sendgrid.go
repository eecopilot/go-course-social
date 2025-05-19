package mailer

import (
	"bytes"
	"errors"
	"text/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewSendGridMailer(apiKey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{fromEmail: fromEmail, apiKey: apiKey, client: client}
}

func (m *SendGridMailer) Send(templateFile string, username, email string, data any, isSandbox bool) error {
	from := mail.NewEmail(FromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// template data
	tmpl, err := template.ParseFS(FS, "/templates/"+templateFile)
	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "Body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "templateFile", body.String())

	// sandbox
	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	// max retry 3 times
	for i := 0; i < MaxRetry; i++ {
		_, err := m.client.Send(message)
		if err != nil {

			// sleep
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		return nil
	}
	return errors.New("failed to send email")
}
