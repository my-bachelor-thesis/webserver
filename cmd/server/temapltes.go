package main

import (
	"errors"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"path/filepath"
	"webserver/internal/config"
)

type Template struct {
	templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return errors.New("Template not found -> " + name)
	}
	return tmpl.ExecuteTemplate(w, name, data)
}

func getAllTemplates() (map[string]*template.Template, error) {
	dir := config.GetInstance().TemplatesDir
	partial, err := filepath.Glob(filepath.Join(dir, "partial", "*.html"))
	if err != nil {
		return nil, err
	}
	all := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		pageName := filepath.Base(page)
		all[pageName] = template.Must(template.ParseFiles(append(partial, page)...))
	}
	return all, nil
}
