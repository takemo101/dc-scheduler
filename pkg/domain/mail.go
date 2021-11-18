package domain

// --- TemplateMailNotify ---

// TemplateData テンプレートで利用するマップデータ
type TemplateData map[string]interface{}

// TemplateMailNotify メール通知情報
type TemplateMailNotify interface {
	ToEmailAddress() string
	Key() string
	Data() TemplateData
}

// --- TemplateMailSender ---

type TemplateMailSender interface {
	Send(notify TemplateMailNotify) error
}
