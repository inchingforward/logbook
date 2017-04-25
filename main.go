package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/echo-contrib/pongor"
	_ "github.com/flosch/pongo2-addons"
	"github.com/gorilla/securecookie"
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

func renderStaticTemplate(c echo.Context, templateName string) error {
	return c.Render(http.StatusOK, templateName, nil)
}

func index(c echo.Context) error {
	return renderStaticTemplate(c, "index.html")
}

func about(c echo.Context) error {
	return renderStaticTemplate(c, "about.html")
}

func getLogin(c echo.Context) error {
	return renderStaticTemplate(c, "login.html")
}

func sessionDump(c echo.Context) error {
	sess, err := store.Get(c.Request(), "session")
	if err != nil {
		log.Printf("%v\n", err)
		return c.Render(http.StatusOK, "error.html", err.Error())
	}

	log.Printf("session  get: %v\n", sess)

	return c.Render(http.StatusBadRequest, "session.html", map[string]interface{}{
		"session": sess,
	})
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Render(http.StatusBadRequest, "login.html", map[string]interface{}{
			"message":  "Username and Password are required.",
			"username": username,
			"password": password,
		})
	}

	log.Printf("login for '%v'", username)
	user, err := models.Login(username, password)
	if err != nil {
		return c.Render(http.StatusBadRequest, "login.html", map[string]interface{}{
			"message":  err.Error(),
			"username": username,
			"password": password,
		})
	}

	sessUser := SessionUser{user.ID, user.UserName}
	sess, _ := store.Get(c.Request(), "session")

	sess.Values["username"] = user.UserName
	sess.Values["userid"] = user.ID

	sess.Save(c.Request(), c.Response())

	log.Printf("session save: %v\n", sess)

	return c.Render(http.StatusOK, "login.html", map[string]interface{}{"message": "Logged in.", "user": sessUser})
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
		return c.Render(http.StatusOK, "error.html", err.Error())
	}
	err = c.Render(http.StatusOK, "user_logbook.html", map[string]interface{}{"logbook": logbook, "username": username, "paginator": pag, "tag": tag})
	if err != nil {
		err = c.Render(http.StatusOK, "error.html", err.Error())
	}

	return err
}

func init() {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=logbook sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	models.SetDB(db)

	store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
}

func main() {
	flag.BoolVar(&debug, "debug", false, "true to enable debug")
	flag.Parse()

	log.Printf("debug: %v\n", debug)

	e := echo.New()

	e.Use(middleware.Logger())

	e.Renderer = pongor.GetRenderer(pongor.PongorOption{
		Reload: debug,
	})

	e.Static("/static", "static")
	e.GET("/", index)
	e.GET("/about", about)
	e.GET("/login", getLogin)
	e.POST("/login", login)
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
	e.GET("/sessiondump", sessionDump)

	e.Logger.Fatal(e.Start(":8006"))
}
