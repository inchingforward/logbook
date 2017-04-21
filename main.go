package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"strconv"

	"github.com/echo-contrib/pongor"
	_ "github.com/flosch/pongo2-addons"
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

type paginator struct {
	entriesPerPage, offset, page, prevPage, nextPage int
}

const (
	entriesPerPage = 20
)

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

func index(c echo.Context) error {
	return renderStaticTemplate(c, "index.html")
}

func about(c echo.Context) error {
	return renderStaticTemplate(c, "about.html")
}

func contact(c echo.Context) error {
	return renderStaticTemplate(c, "contact.html")
}

func makePaginator(c echo.Context) paginator {
	pageParam := c.QueryParam("page")
	page, error := strconv.Atoi(pageParam)

	if error != nil {
		page = 1
	}

	nextPage := page + 1
	prevPage := page - 1

	offset := (page - 1) * entriesPerPage

	if prevPage < 0 {
		prevPage = 0
	}

	return paginator{entriesPerPage, offset, page, prevPage, nextPage}
}

func getUserLogbook(c echo.Context) error {
	username := c.Param("username")
	pag := makePaginator(c)

	logbook, err := models.GetUserLogbook(username, pag.offset, pag.entriesPerPage)
	if err != nil {
		log.Printf("Error getting user logbook: %v\n", err)
		return c.Render(http.StatusOK, "error.html", err.Error())
	}
	err = c.Render(http.StatusOK, "user_logbook.html", map[string]interface{}{"logbook": logbook, "username": username, "paginator": pag})
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

	/*t := &Template{
		templates: template.Must(template.New("main").Funcs(funcMap).ParseGlob("templates/*.html")),
	}*/
	e.Renderer = pongor.GetRenderer(pongor.PongorOption{
		Reload: debug,
	})

	//e.Renderer = t
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
