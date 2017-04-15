package main

import (
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

var (
	debug = false
)

func notYetImplemented(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}

func renderStaticTemplate(c echo.Context, templateName string) error {
	return c.Render(http.StatusOK, templateName, nil)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if debug {
		t.templates = template.Must(template.ParseGlob("templates/*.html"))
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func index(c echo.Context) error {
	return renderStaticTemplate(c, "index.html")
}

func about(c echo.Context) error {
	return renderStaticTemplate(c, "about.html")
}

func contact(c echo.Context) error {
	return renderStaticTemplate(c, "contact.html")
}

func main() {
	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	log.Printf("debug: %v\n", debug)

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Renderer = t
	e.Static("/static", "static")
	e.GET("/", index)
	e.GET("/about", about)
	e.GET("/contact", contact)
	e.GET("/login", notYetImplemented)
	e.POST("/login", notYetImplemented)
	e.POST("/logout", notYetImplemented)
	e.GET("/logbook", notYetImplemented)
	e.GET("/logbook/add", notYetImplemented)
	e.POST("/logbook/add", notYetImplemented)
	e.GET("/logbook/:uuid", notYetImplemented)
	e.POST("/logbook/:uuid", notYetImplemented)
	e.POST("/fetchtitle", notYetImplemented)
	e.GET("/users/:username", notYetImplemented)
	e.GET("/users/:username/logbook", notYetImplemented)
	e.GET("/users/:username/logbook/:uuid", notYetImplemented)

	e.Logger.Fatal(e.Start(":8000"))
}
