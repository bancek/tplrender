package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func main() {
	var templatePath string
	var destination string
	var leftDelim string
	var rightDelim string

	flag.StringVar(&templatePath, "template", "-", "The template to render")
	flag.StringVar(&destination, "dest", "-", "Filename for the rendered template (default stdout)")
	flag.StringVar(&leftDelim, "leftdelim", "{{", "Left-hand side delimiter for the template")
	flag.StringVar(&rightDelim, "rightdelim", "}}", "Right-hand side delimiter for the template")

	flag.Parse()

	var input io.Reader = os.Stdin

	if templatePath != "-" {
		f, err := os.Open(templatePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "faild to open templatePath file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	}

	templateString, err := ioutil.ReadAll(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read template: %v\n", err)
		os.Exit(1)
	}

	tmpl, err := template.New("tplrender").Funcs(sprig.TxtFuncMap()).Delims(leftDelim, rightDelim).Parse(string(templateString))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse template: %v\n", err)
		os.Exit(1)
	}

	env := map[string]string{}

	for _, item := range os.Environ() {
		parts := strings.SplitN(item, "=", 2)
		env[parts[0]] = parts[1]
	}

	data := struct {
		Args []string
		Env  map[string]string
	}{
		Args: flag.Args(),
		Env:  env,
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Fprintf(os.Stderr, "failed to render template: %v\n", err)
		os.Exit(1)
	}

	var output io.Writer = os.Stdout

	if destination != "-" {
		f, err := os.Create(destination)
		if err != nil {
			fmt.Fprintf(os.Stderr, "faild to create destination file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		output = f
	}

	if _, err := io.Copy(output, &buf); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to the destination file: %v\n", err)
		os.Exit(1)
	}
}
