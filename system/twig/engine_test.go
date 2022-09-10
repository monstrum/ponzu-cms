package twig

import (
	"bytes"
	"github.com/monstrum/ponzu-cms/system/item"
	"testing"
)

func TestName(t *testing.T) {
	c, err := Twig.Render("test.html.twig", map[string]interface{}{"Key": "Hello"})
	if err != nil {
		t.Errorf("error is not supposed to happened")
	}
	buf := bytes.NewBufferString("test data Hello")
	if !bytes.Equal(c, buf.Bytes()) {
		t.Errorf("test data is wrong! \"%s\" and \"%s\"", buf.Bytes(), c)
	}
	tes, err := Twig.Render("admin.html.twig", map[string]interface{}{
		"Logo":    "cfg",
		"Types":   item.Types,
		"Subview": "we",
	})
	_ = tes
	if err != nil {
		t.Errorf("error is not supposed to happened \"%s\"", err)
	}
}
