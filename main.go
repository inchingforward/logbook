package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/inchingforward/logbook/models"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/russross/blackfriday"
)

type Template struct {
	templates *template.Template
}

var (
	debug   = false
	funcMap template.FuncMap
)

func notYetImplemented(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}

func renderStaticTemplate(c echo.Context, templateName string) error {
	return c.Render(http.StatusOK, templateName, nil)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if debug {
		t.templates = template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html"))
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

func getUserLogbook(c echo.Context) error {
	username := c.Param("username")

	logbook, err := models.GetUserLogbook(username)
	if err != nil {
		log.Printf("Error getting user logbook: %v\n", err)
		return c.Render(http.StatusOK, "error.html", err.Error())
	}

	err = c.Render(http.StatusOK, "user_logbook.html", struct {
		Logbook []*models.Entry
		User    string
	}{logbook, username})
	if err != nil {
		err = c.Render(http.StatusOK, "error.html", err.Error())
	}

	return err
}

func markDownBasic(args ...interface{}) template.HTML {
	s := blackfriday.MarkdownBasic([]byte(fmt.Sprintf("%s", args...)))
	return template.HTML(s)
}

func init() {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=logbook sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	models.SetDB(db)

	funcMap = template.FuncMap{
		"mdb": markDownBasic,
	}
}

func main() {
	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	log.Printf("debug: %v\n", debug)

	e := echo.New()

	e.Use(middleware.Logger())

	t := &Template{
		templates: template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html")),
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
	e.GET("/users/:username/logbook", getUserLogbook)
	e.GET("/users/:username/logbook/:uuid", notYetImplemented)

	e.Logger.Fatal(e.Start(":8006"))
}
