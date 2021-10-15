package core

import (
	"fmt"
	"strings"

	"github.com/flosch/pongo2/v4"
	"github.com/gofiber/template/django"
)

type BindData map[string]interface{}

// TemplateEngine is view engine
type TemplateEngine struct {
	Engine *django.Engine
}

// NewTemplateEngine set up template
func NewTemplateEngine(
	config Config,
	p Path,
	logger Logger,
) TemplateEngine {
	engine := django.New(config.Template.Path, config.Template.Suffix)
	engine.Reload(config.Template.Reload)
	engine.Debug(config.App.Debug)

	engine.AddFunc("nl2br", func(value interface{}) string {
		if str, ok := value.(string); ok {
			return strings.Replace(str, "\n", "<br />", -1)
		}
		return ""
	})

	engine.AddFunc("static", func(value interface{}) string {
		return HandleTemplateFuncURL(value, func(str string) string {
			return p.StaticURL(str)
		})
	})

	engine.AddFunc("public", func(value interface{}) string {
		return HandleTemplateFuncURL(value, func(str string) string {
			return p.PublicURL(str)
		})
	})

	engine.AddFunc("url", func(value interface{}, a ...interface{}) string {
		if str, ok := value.(string); ok {
			return p.URL(str, a...)
		}
		return ""
	})

	return TemplateEngine{
		Engine: engine,
	}
}

// ParseTemplate parse from template file
func (template TemplateEngine) ParseTemplate(path string, data BindData) (string, error) {
	tmpl, ok := template.Engine.Templates[path]
	if !ok {
		return "", fmt.Errorf("template %s does not exist", path)
	}

	return tmpl.Execute(pongo2.Context(data))
}

// ParseFromString parse from string
func (template TemplateEngine) ParseFromString(text string, data BindData) (string, error) {
	tmpl, err := pongo2.FromString(text)
	if err != nil {
		return "", err
	}
	return tmpl.Execute(pongo2.Context(data))
}

// HandleTemplateFuncURL template func
func HandleTemplateFuncURL(value interface{}, handle func(string) string) string {
	if str, ok := value.(string); ok {
		for _, prefix := range []string{"https://", "http://"} {
			if strings.Contains(str, prefix) {
				return str
			}
		}
		return handle(str)
	}
	return ""
}
