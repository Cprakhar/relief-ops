package mail

import (
	"embed"

	"github.com/cprakhar/relief-ops/services/user-service/repo"
)

const (
	FromName            = "Relief Ops"
	MaxRetries          = 3
	AdminNotifyTemplate = "admin_notify.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) (int, error)
	SpamMail(users []*repo.User, data any, isSandbox bool) error
}
