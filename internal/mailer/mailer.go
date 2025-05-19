package mailer

import (
	"embed"
)

// hard code
const (
	FromName               = "EECO PILOT"
	MaxRetry               = 3
	UserInvitationTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile string, username, email string, data any, isSandbox bool) error
}
