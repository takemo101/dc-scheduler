package infrastructure

import (
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- TemplateMailSender ---

// TemplateMailSender テンプレートメール
type TemplateMailSender struct {
	mailTemplate helper.MailTemplate
}

// NewTemplateMailSender コンストラクタ
func NewTemplateMailSender(
	mailTemplate helper.MailTemplate,
) domain.TemplateMailSender {
	return TemplateMailSender{
		mailTemplate,
	}
}

// Send テンプレートメールでの送信処理
func (sender TemplateMailSender) Send(notify domain.TemplateMailNotify) (err error) {
	mail, _, err := sender.mailTemplate.GetTextMailCreatorByKey(
		notify.Key(),
		core.BindData(notify.Data()),
	)
	if err != nil {
		return err
	}

	mail.To(notify.ToEmailAddress())

	return mail.Send()
}
