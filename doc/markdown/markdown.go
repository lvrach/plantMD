package markdown

import (
	"fmt"
	"strings"
)

type Element string

func Image(name, path string) Element {
	return Element(fmt.Sprintf("![%s](%s)", name, path))
}

func H1(title string) Element {
	return Element(fmt.Sprintf("\n# %s\n", title))
}

func H2(title string) Element {
	return Element(fmt.Sprintf("\n## %s\n", title))
}

func Text(text string) Element {
	return Element(text)
}

type Document struct {
	buffer string
}

func (d *Document) Append(e Element) {
	d.buffer = fmt.Sprintf("%s\n%s", d.buffer, string(e))
}

func (d *Document) Reader() *strings.Reader {
	return strings.NewReader(d.buffer)
}
