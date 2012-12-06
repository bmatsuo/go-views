package main

import (
	"fmt"
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
var outLogger = log.New(os.Stdout, "go-views", log.LstdFlags)
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

func findDirectories(root string) ([]string, error) {
	dirs := make([]string, 0, 10)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func isTemplate(path string) bool {
	return strings.HasSuffix(path, opt.Extension)
}

func findTemplates(root string) ([]string, error) {
	paths := make([]string, 0, 10)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if isTemplate(path) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func templateMap(root string, paths []string) (map[string]string, error) {
	m := make(map[string]string, len(paths))
	for _, path := range paths {
		relPath := path[len(root)+1:]
		name := relPath[:len(relPath)-len(opt.Extension)]
		body, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if len(body) > 0 {
			body = body[:len(body)-1]
		}
		m[name] = string(body)
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

func compile(root string) error {
	paths, err := findTemplates(root)
	if err != nil {
		return fmt.Errorf("error locating templates: %v", err)
	}

	m, err := templateMap(root, paths)
	if err != nil {
		return fmt.Errorf("error reading files: %v", err)
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
		return fmt.Errorf("error generating view file: %v", err)
	}
	return nil
}

func watchAndRecompile(root string) error {
	dirs, err := findDirectories(root)
	if err != nil {
		errLogger.Fatalln("discovery error: ", err)
	}

	done := make(chan error)
	recompile := make(chan string, 1)
	go func() {
		for path := range recompile {
			outLogger.Println("recompile triggered for ", path)
			err := compile(root)
			if err != nil {
				errLogger.Println(err)
			}
		}
		done <- nil
	}()

	watcher, err := watch(recompile, isTemplate)
	if err != nil {
		close(recompile)
		errLogger.Fatalln("error initializing watcher: ", err)
	}

	for _, d := range dirs {
		err := watcher.Watch(d)
		if err != nil {
			watcher.Close()
			errLogger.Fatalln("watcher populate error: ", err)
		}
	}

	return <-done
}

func init() {
	opt, args = parseOptions()
}

func main() {
	root, err := filepath.Abs(args[0])
	if err != nil {
		errLogger.Fatalln("error locating root: ", err)
	}

	err = compile(root)
	if err != nil {
		// probably shouldn't be fatal with opt.Watch
		errLogger.Fatalln(err)
	}

	if opt.Watch {
		err = watchAndRecompile(root)
		if err != nil {
			errLogger.Fatalln(err)
		}
	}
}
