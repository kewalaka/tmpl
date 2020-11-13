package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	htmltemplate "html/template"
	"text/template"

	"github.com/Masterminds/sprig"
)

var (
	flagInput = flag.String("f", "-", "Input source")
	flagHTML  = flag.Bool("html", false, "If true, use html/tmplate instead of text/template")
)

func main() {
	flag.Parse()
	in, err := getInput(*flagInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening input:", err)
		os.Exit(1)
	}
	if err := tmpl(in, *flagHTML, os.Stdout, envMap()); err != nil {
		fmt.Fprintln(os.Stderr, "error opening input:", err)
		os.Exit(1)
	}
}

func getInput(path string) (io.Reader, error) {
	if path == "-" {
		return os.Stdin, nil
	}
	return os.Open(path)
}

func envMap() map[string]string {
	result := map[string]string{}
	for _, envvar := range os.Environ() {
		parts := strings.SplitN(envvar, "=", 2)
		result[parts[0]] = parts[1]
	}
	return result
}

func tmpl(in io.Reader, htmlMode bool, out io.Writer, ctx interface{}) error {
	i, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}

	if htmlMode {
		tmpl, err := htmltemplate.New("format string").Funcs(sprig.HtmlFuncMap()).Parse(string(i))
		if err != nil {
			return err
		}
		return tmpl.Execute(out, ctx)
	} else {
		tmpl, err := template.New("format string").Funcs(sprig.TxtFuncMap()).Parse(string(i))
		if err != nil {
			return err
		}
		return tmpl.Execute(out, ctx)
	}
}
