package core

import (
	"fmt"
	"path"
	"strings"
)

// path is path helper
type Path struct {
	config Config
}

func NewPath(
	config Config,
) Path {
	return Path{
		config: config,
	}
}

func (p Path) URL(suffix string, a ...interface{}) string {
	return p.config.App.URL + "/" + fmt.Sprintf(suffix, a...)
}

func (p Path) StaticURL(suffix string) string {
	return p.URL(path.Join(p.config.Static.Prefix, suffix))
}

func (p Path) PublicURL(suffix string) string {
	return p.StaticURL(path.Join(p.config.File.Public, suffix))
}

func (p Path) Current(suffix string) string {
	return path.Join(p.config.File.Current, suffix)
}

func (p Path) Storage(suffix string) string {
	return p.Current(path.Join(p.config.File.Storage, suffix))
}

func (p Path) Static(suffix string) string {
	return p.Current(path.Join(p.config.Static.Root, suffix))
}

func (p Path) Public(suffix string) string {
	return p.Storage(path.Join(p.config.File.Public, suffix))
}

func (p Path) StorageSubstract(str string) string {
	storage := p.Storage("")
	return strings.Replace(str, storage, "", 1)
}

func (p Path) StaticSubstract(str string) string {
	static := p.Static("")
	return strings.Replace(str, static, "", 1)
}

func (p Path) PublicSubstract(str string) string {
	public := p.Public("")
	return strings.Replace(str, public, "", 1)
}
