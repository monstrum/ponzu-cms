package twig

import (
	"embed"
	"github.com/tyler-sommer/stick"
	"path/filepath"
)

//go:embed templates
var templates embed.FS

type twigLoader struct {
}

func (l *twigLoader) Load(name string) (stick.Template, error) {
	b, err := templates.ReadFile(filepath.Join("templates", name))
	if err != nil {
		return nil, err
	}
	return &byteTemplate{b}, nil
}
