package core

import (
	"io"
	"log"
	"net/smtp"
	"strconv"

	"github.com/domodwyer/mailyak/v3"
)

// MailCreator interface
type MailCreator interface {
	From(string)
	FromName(string)
	To(...string)
	Bcc(...string)
	Subject(string)
	ReplyTo(string)
	Attach(string, io.Reader)
	Text(string)
	HTML(string)
	Parse(string, BindData) (string, error)
	TemplateText(string, BindData) string
	TemplateHTML(string, BindData) string
	Send() error
}

// MailFactory interface
type MailFactory interface {
	Create() MailCreator
}

// YakMailCreator mailyak creator struct
type YakMailCreator struct {
	mail     *mailyak.MailYak
	template TemplateEngine
	logger   Logger
}

// YakMailFactory mailyak factory struct
type YakMailFactory struct {
	config   SMTP
	template TemplateEngine
	logger   Logger
}

// NewMailFactory constructor
func NewMailFactory(
	config Config,
	template TemplateEngine,
	logger Logger,
) MailFactory {

	// template is once load
	if err := template.Engine.Load(); err != nil {
		logger.Error(err)
	}

	return YakMailFactory{
		config:   config.SMTP,
		template: template,
		logger:   logger,
	}
}

func (creator *YakMailCreator) From(address string) {
	creator.mail.From(address)
}

func (creator *YakMailCreator) FromName(name string) {
	creator.mail.FromName(name)
}

func (creator *YakMailCreator) To(to ...string) {
	creator.mail.To(to...)
}

func (creator *YakMailCreator) Bcc(bcc ...string) {
	creator.mail.Bcc(bcc...)
}

func (creator *YakMailCreator) Subject(subject string) {
	creator.mail.Subject(subject)
}

func (creator *YakMailCreator) ReplyTo(reply string) {
	creator.mail.ReplyTo(reply)
}

func (creator *YakMailCreator) Attach(name string, reader io.Reader) {
	creator.mail.Attach(name, reader)
}

func (creator *YakMailCreator) Text(text string) {
	creator.mail.Plain().Set(text)
}

func (creator *YakMailCreator) HTML(html string) {
	creator.mail.HTML().Set(html)
}

func (creator *YakMailCreator) Parse(text string, data BindData) (string, error) {
	return creator.template.ParseFromString(text, data)
}

func (creator *YakMailCreator) TemplateText(path string, data BindData) string {
	parse, err := creator.template.ParseTemplate(path, data)
	if err != nil {
		creator.logger.Error(err)
	} else {
		creator.Text(parse)
	}

	return parse
}

func (creator *YakMailCreator) TemplateHTML(path string, data BindData) string {
	parse, err := creator.template.ParseTemplate(path, data)
	if err != nil {
		creator.logger.Error(err)
	} else {
		creator.HTML(parse)
	}

	return parse
}

func (creator *YakMailCreator) Send() error {
	return creator.mail.Send()
}

// Create mail creator
func (factory YakMailFactory) Create() MailCreator {
	auth := smtp.PlainAuth(
		factory.config.Identity,
		factory.config.User,
		factory.config.Pass,
		factory.config.Host,
	)
	var mail *mailyak.MailYak
	host := factory.config.Host + ":" + strconv.Itoa(factory.config.Port)

	if factory.config.Encryption == "tls" {
		var err error

		mail, err = mailyak.NewWithTLS(
			host,
			auth,
			nil,
		)
		if err != nil {
			log.Fatalf("fail to yakmail : %v", err)
		}
	} else {

		mail = mailyak.New(
			host,
			auth,
		)
	}
	mail.FromName(factory.config.From.Name)
	mail.From(factory.config.From.Address)

	return &YakMailCreator{
		mail:     mail,
		template: factory.template,
		logger:   factory.logger,
	}
}
