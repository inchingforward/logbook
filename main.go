package main

import (
	"encoding/gob"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/flosch/pongo2"
	_ "github.com/flosch/pongo2-addons"
	"github.com/gorilla/sessions"
	"github.com/inchingforward/logbook/models"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

type paginator struct {
	entriesPerPage, offset, page, prevPage, nextPage int
}

// Renderer renders templates.
type Renderer struct {
	TemplateDir string
	Reload      bool
}

const (
	entriesPerPage = 20
)

var (
	debug = false
	store sessions.Store
)

// A SessionUser is stored in a user's session.
type SessionUser struct {
	ID       uint64
	UserName string
}

func notYetImplemented(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}

func renderTemplate(c echo.Context, templateName string) error {
	return c.Render(http.StatusOK, templateName, pongo2.Context{})
}

func renderError(c echo.Context, message string) error {
	return c.Render(http.StatusOK, "error.html", pongo2.Context{"error": message})
}

func index(c echo.Context) error {
	return renderTemplate(c, "index.html")
}

func about(c echo.Context) error {
	return renderTemplate(c, "about.html")
}

func getLogin(c echo.Context) error {
	return renderTemplate(c, "login.html")
}

func sessionDump(c echo.Context) error {
	return renderTemplate(c, "session.html")
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Render(http.StatusBadRequest, "login.html", pongo2.Context{
			"message":  "Username and Password are required.",
			"username": username,
			"password": password,
		})
	}

	log.Printf("login for '%v'", username)

	user, err := models.Login(username, password)
	if err != nil {
		return c.Render(http.StatusBadRequest, "login.html", pongo2.Context{
			"message":  err.Error(),
			"username": username,
			"password": password,
		})
	}

	sessUser := SessionUser{user.ID, user.UserName}
	sess, err := store.Get(c.Request(), "session")
	if err != nil {
		log.Printf("%v\n", err)
		return c.Redirect(http.StatusFound, "/login")
	}

	sess.Values["User"] = sessUser
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/sessiondump")
}

func logout(c echo.Context) error {
	sess, _ := store.Get(c.Request(), "session")

	sess.Values["User"] = nil
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/login")
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
	tag := c.QueryParam("tag")

	logbook, err := models.GetUserLogbook(username, tag, pag.offset, pag.entriesPerPage)
	if err != nil {
		log.Printf("Error getting user logbook: %v\n", err)
		return renderError(c, err.Error())
	}

	err = c.Render(http.StatusOK, "user_logbook.html", pongo2.Context{"logbook": logbook, "username": username, "paginator": pag, "tag": tag})
	if err != nil {
		err = c.Render(http.StatusOK, "error.html", pongo2.Context{"error": err.Error()})
	}

	return err
}

// Render renders a pongo2 template.
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	filename := path.Join(r.TemplateDir, name)
	template := pongo2.Must(pongo2.FromFile(filename))

	pctx := data.(pongo2.Context)

	sess, err := store.Get(c.Request(), "session")
	if err == nil {
		pctx["session"] = sess
	}

	return template.ExecuteWriter(pctx, w)
}

func init() {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=logbook sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	models.SetDB(db)

	storeKey := os.Getenv("LOGBOOK_STORE_KEY")
	if storeKey == "" {
		log.Fatal("The LOGBOOK_STORE_KEY is not set")
	}
	store = sessions.NewCookieStore([]byte(storeKey))

	gob.Register(SessionUser{})
}

func main() {
	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	log.Printf("debug: %v\n", debug)

	e := echo.New()

	e.Use(middleware.Logger())

	e.Renderer = &Renderer{TemplateDir: "templates", Reload: debug}

	e.Static("/static", "static")
	e.GET("/", index)
	e.GET("/about", about)
	e.GET("/login", getLogin)
	e.POST("/login", login)
	e.POST("/logout", logout)
	e.GET("/logbook", notYetImplemented)
	e.GET("/logbook/add", notYetImplemented)
	e.POST("/logbook/add", notYetImplemented)
	e.GET("/logbook/:uuid", notYetImplemented)
	e.POST("/logbook/:uuid", notYetImplemented)
	e.POST("/fetchtitle", notYetImplemented)
	e.GET("/users/:username", notYetImplemented)
	e.GET("/users/:username/logbook", getUserLogbook)
	e.GET("/users/:username/logbook/:uuid", notYetImplemented)
	e.GET("/sessiondump", sessionDump)

	e.Logger.Fatal(e.Start(":8006"))
}
