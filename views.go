package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var opt *Options
var args []string
var errLogger = log.New(os.Stderr, "go-views", log.LstdFlags)
var templatePaths []string
var viewFile = "views.go"
var viewTemplate = `// This file is auto-generated.
// Do not commit it to any version control system.
package {|or .package "views"|}

import (
	"io"
	"text/template"
)

var Templates = template.Must(template.New("views").Parse(` + "`" + `
{|range $name, $template := .templates|}{{define "{|$name|}"}}{|$template|}{{end}}
{|end|}` + "`" + `))

func ExecuteTemplate(w io.Writer, name string, data interface{}) error {
	return Templates.ExecuteTemplate(w, name, data)
}
func Render(w io.Writer, name string, data interface{}) error {
	return ExecuteTemplate(w, name, data)
}
`

func findTemplates(root string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, opt.Extension) {
			templatePaths = append(templatePaths, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func templateMap(root string) (map[string]string, error) {
	m := make(map[string]string, len(templatePaths))
	for _, path := range templatePaths {
		relPath := path[len(root)+1:]
		name := relPath[:len(relPath)-len(opt.Extension)]
		body, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		m[name] = string(body[:len(body)-1])
	}
	return m, nil
}

func makeViewFile(w io.Writer, data interface{}) error {
	t, err := template.New("viewTemplate").Delims("{|", "|}").Parse(viewTemplate)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

func compile(root string) {
	err := findTemplates(root)
	if err != nil {
		errLogger.Fatalln("error locating templates: ", err)
	}

	m, err := templateMap(root)
	if err != nil {
		errLogger.Fatalln("error reading files: ", err)
	}
	data := map[string]interface{}{
		"package":   opt.Package,
		"templates": m,
	}

	out := os.Stdout
	if opt.Output != "-" {
		out, err = os.OpenFile(opt.Output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			errLogger.Println("error opening output: ", err)
		}
		defer out.Close()
	}
	err = makeViewFile(out, data)
	if err != nil {
		errLogger.Fatalln("error generating view file: ", err)
	}
}

func init() {
	opt, args = parseOptions()
}

func main() {
	root, err := filepath.Abs(args[0])
	if err != nil {
		errLogger.Fatalln("error locating root: ", err)
	}
	compile(root)
}
