package twig

type twigWriter struct {
	content []byte
}

func (w *twigWriter) Write(p []byte) (int, error) {
	w.content = append(w.content[:], p[:]...)
	return len(w.content), nil
}

func (w *twigWriter) getContent() []byte {
	c := w.content
	w.content = []byte{}
	return c
}
