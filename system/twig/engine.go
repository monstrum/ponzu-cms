package twig

import (
	"github.com/tyler-sommer/stick"
)

var (
	Twig = createTwigEngine()
)

// Engine Twig engine interface
type Engine interface {
	Render(template string, ctx map[string]interface{}) (content []byte, err error)
}

func createTwigEngine() Engine {
	engine := twigEngine{
		stick.New(&twigLoader{}),
		&twigWriter{},
	}
	return &engine
}

type twigEngine struct {
	renderer *stick.Env
	writer   *twigWriter
}

func (t *twigEngine) Render(template string, ctx map[string]interface{}) ([]byte, error) {
	v := make(map[string]stick.Value)
	for i, c := range ctx {
		v[i] = c
	}
	err := t.renderer.Execute(template, t.writer, v)
	if err != nil {
		return []byte{}, err
	}
	return t.writer.getContent(), nil
}
