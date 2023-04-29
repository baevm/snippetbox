package templates

import (
	"io/fs"
	"path/filepath"
	"snippetbox/internal/models"
	"snippetbox/ui"
	"text/template"
	"time"
)

type TemplateData struct {
	Account         *models.User
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func HumanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": HumanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")

	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)

		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

// func (app *App) Render(w http.ResponseWriter, status int, page string, data *TemplateData) {
// 	ts, ok := app.templateCache[page]

// 	if !ok {
// 		app.serverError(w, fmt.Errorf("template %s doesnt exist", page))
// 		return
// 	}

// 	buf := new(bytes.Buffer)

// 	err := ts.ExecuteTemplate(buf, "base", data)

// 	if err != nil {
// 		app.serverError(w, err)
// 		return
// 	}

// 	w.WriteHeader(status)
// 	buf.WriteTo(w)
// }
