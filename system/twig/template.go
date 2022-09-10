package twig

import (
	"bytes"
	"io"
	"io/fs"
)

type strTemplate struct {
	content string
}

func (t *strTemplate) Name() string {
	return ""
}

func (t *strTemplate) Contents() io.Reader {
	return bytes.NewBufferString(t.content)
}

type byteTemplate struct {
	content []byte
}

func (b *byteTemplate) Name() string {
	return ""
}

func (b *byteTemplate) Contents() io.Reader {
	return bytes.NewReader(b.content)
}

type fileTemplate struct {
	file fs.File
}

func (f *fileTemplate) Name() string {
	return ""
}

func (f *fileTemplate) Contents() io.Reader {
	var content []byte
	var a, err = f.file.Read(content)
	if err != nil {
		panic(err)
	}
	if a == 0 {
		panic(a)
	}
	return bytes.NewReader(content)
}
