package view

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

var (
	t     *Template
	debug = false
)

// Template represents the parsed templates from the "templates" directory.
type Template struct {
	templates *template.Template
}

func init() {
	t = &Template{
		templates: template.Must(template.New("main").ParseGlob("templates/*.html")),
	}
}

// SetRenderer sets the Logbook templates as the Echo engine's renderer.
func SetRenderer(e *echo.Echo, debugEnabled bool) {
	debug = debugEnabled
	e.SetRenderer(t)
}

// Render renders the template referenced by name and passes the data value
// into the template.  If main was run with the debug argument, the templates
// are re-parsed on each call to Render.
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if debug {
		t.templates = template.Must(template.New("main").ParseGlob("templates/*.html"))
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

// RenderTemplate renders a simple template with no data.
func RenderTemplate(c echo.Context, templ string) error {
	return c.Render(http.StatusOK, templ, nil)
}
