package main

import (
	"io"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/lvrach/plantMD/doc/markdown"
	"github.com/lvrach/plantMD/render/puml"
)

type UMLProcessor struct {
	client puml.Client
}

func (p UMLProcessor) process(file string) *os.File {
	uml, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	img, err := p.client.Render(uml)
	if err != nil {
		log.Fatal(err)
	}
	out, err := os.Create(strings.ReplaceAll(file, "puml", "png"))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, img)
	if err != nil {
		log.Fatal(err)
	}

	return out
}

func main() {
	p := UMLProcessor{
		client: puml.Client{
			Host: "http://localhost:8080",
		},
	}

	d := markdown.Document{}
	d.Append(markdown.H1("UML diagrams"))

	wg := sync.WaitGroup{}
	images := make(chan string, 10)

	go func() {

	}()

	wg.Add(len(os.Args[1:]))
	for _, umlFile := range os.Args[1:] {
		go func(umlFile string) {
			images <- p.process(umlFile).Name()
			wg.Done()
		}(umlFile)
	}

	go func() {
		wg.Wait()
		close(images)
	}()

	for imgName := range images {
		log.Println(imgName)
		d.Append(markdown.H2(pathToTitle(imgName)))
		d.Append(markdown.Image("uml", "./"+path.Clean(imgName)))
	}
	out, err := os.Create("./uml.md")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, d.Reader())
	if err != nil {
		log.Fatal(err)
	}
}

func pathToTitle(p string) string {
	p = path.Clean(p)
	p = strings.ReplaceAll(p, "/", " ")
	p = strings.TrimSuffix(p, path.Ext(p))
	return p
}
