package main

import (
	"fmt"
	"github.com/russross/blackfriday"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func Run(input []byte, w io.Writer) {
	r := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		Flags: blackfriday.CommonHTMLFlags | blackfriday.CompletePage,
	})
	p := blackfriday.New(
		blackfriday.WithRenderer(r),
		blackfriday.WithExtensions(blackfriday.CommonExtensions),
	)
	ast := p.Parse(input)
	r.RenderHeader(w, ast)
	fmt.Fprintf(w, "<style>%s</style>", GFMCSS)
	fmt.Fprintf(w, `<div class="markdown-body">`)
	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.RenderNode(w, node, entering)
	})
	fmt.Fprintf(w, `</div>`)
	r.RenderFooter(w, ast)
}

type SourceLoader interface {
	Load(string, io.Writer) error
}

type FileLoader string

func (l *FileLoader) Load(path string, dest io.Writer) error {
	content, err := ioutil.ReadFile(string(*l))
	if err != nil {
		return fmt.Errorf("failed to load file %s", *l)
	} else {
		Run(content, dest)
	}
	return nil
}

type OnceOffLoader []byte

func (l *OnceOffLoader) Load(path string, dest io.Writer) error {
	Run(*l, dest)
	return nil
}

type URLLoader string

func (l *URLLoader) Load(path string, dest io.Writer) error {
	resp, err := http.DefaultClient.Get(string(*l))
	if err != nil {
		return fmt.Errorf("failed to load url %s", string(*l))
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to load file %s", *l)
	} else {
		Run(content, dest)
	}
	return nil
}

type DirectoryLoader string

func (l *DirectoryLoader) Load(path string, dest io.Writer) error {
	if strings.Contains(path, "..") {
		return fmt.Errorf("cannot load file with dotted section")
	}

	fp := filepath.Join(string(*l), path)
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return fmt.Errorf("failed to load file %s", fp)
	} else {
		Run(content, dest)
	}
	return nil
}
