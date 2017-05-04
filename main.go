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
	"strings"

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
	TemplateDir   string
	Reload        bool
	TemplateCache map[string]*pongo2.Template
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
		return logout(c)
	}

	sess.Values["User"] = sessUser
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/logbook")
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

	logbook, err := models.GetUserPublicLogbook(username, tag, pag.offset, pag.entriesPerPage)
	if err != nil {
		log.Printf("Error getting user logbook: %v\n", err)
		return renderError(c, err.Error())
	}

	err = c.Render(http.StatusOK, "user_logbook.html", pongo2.Context{
		"logbook":   logbook,
		"username":  username,
		"paginator": pag,
		"tag":       tag,
	})

	if err != nil {
		err = c.Render(http.StatusOK, "error.html", pongo2.Context{"error": err.Error()})
	}

	return err
}

func getLogbook(c echo.Context) error {
	user := getUser(c)
	pag := makePaginator(c)
	tag := c.QueryParam("tag")

	logbook, err := models.GetLogbook(user.ID, tag, pag.offset, pag.entriesPerPage)
	if err != nil {
		log.Printf("Error getting user logbook: %v\n", err)
		return renderError(c, err.Error())
	}

	err = c.Render(http.StatusOK, "logbook.html", pongo2.Context{
		"logbook":   logbook,
		"paginator": pag,
		"tag":       tag,
	})

	if err != nil {
		err = c.Render(http.StatusOK, "error.html", pongo2.Context{"error": err.Error()})
	}

	return err
}

func getAddEntry(c echo.Context) error {
	return c.Render(http.StatusOK, "logbook_entry_add.html", pongo2.Context{})
}

func addEntry(c echo.Context) error {
	entry, err := getFormEntry(c)

	user := getUser(c)
	err = models.InsertEntry(user.ID, entry)
	if err != nil {
		return renderError(c, err.Error())
	}

	return c.Redirect(http.StatusFound, "/logbook")
}

func updateEntry(c echo.Context) error {
	user := getUser(c)
	uuid := c.FormValue("uuid")

	entry, err := models.GetLogbookEntry(user.ID, uuid)
	if err != nil || entry.ID == 0 {
		return renderError(c, "Invalid entry.")
	}

	formEntry, err := getFormEntry(c)
	if err != nil {
		return renderError(c, err.Error())
	}

	entry.Title = formEntry.Title
	entry.URL = formEntry.URL
	entry.Notes = formEntry.Notes
	entry.Private = formEntry.Private
	entry.Tags = formEntry.Tags

	err = models.UpdateEntry(&entry)
	if err != nil {
		return renderError(c, err.Error())
	}

	return c.Redirect(http.StatusFound, "/logbook")
}

func getFormEntry(c echo.Context) (*models.Entry, error) {
	entry := new(models.Entry)
	if err := c.Bind(entry); err != nil {
		return nil, err
	}

	tagStr := c.FormValue("tags")
	tags := strings.Split(tagStr, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	entry.Tags = tags

	return entry, nil
}

func getEntry(c echo.Context) error {
	user := getUser(c)
	entryUUID := c.Param("uuid")

	if entryUUID == "" {
		return c.Render(http.StatusNotFound, "404.html", pongo2.Context{})
	}

	entry, err := models.GetLogbookEntry(user.ID, entryUUID)
	if err != nil {
		return c.Render(http.StatusOK, "message.html", pongo2.Context{
			"message": err.Error(),
		})
	}

	return c.Render(http.StatusOK, "logbook_entry_edit.html", pongo2.Context{
		"entry": entry,
	})
}

func getUser(c echo.Context) SessionUser {
	sess, _ := store.Get(c.Request(), "session")
	return sess.Values["User"].(SessionUser)
}

// GetTemplate returns a template, loading it every time if reload is true.
func (r *Renderer) GetTemplate(name string, reload bool) *pongo2.Template {
	filename := path.Join(r.TemplateDir, name)

	if r.Reload {
		return pongo2.Must(pongo2.FromFile(filename))
	}

	return pongo2.Must(pongo2.FromCache(filename))
}

// Render renders a pongo2 template.
func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	template := r.GetTemplate(name, debug)
	pctx := data.(pongo2.Context)

	sess, err := store.Get(c.Request(), "session")
	if err == nil {
		pctx["session"] = sess
	}

	pctx["csrf"] = c.Get("csrf")

	return template.ExecuteWriter(pctx, w)
}

func ensureSessionUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := store.Get(c.Request(), "session")
		if err != nil {
			return logout(c)
		}

		_, found := sess.Values["User"].(SessionUser)
		if !found {
			return logout(c)
		}

		return next(c)
	}
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
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf",
	}))

	e.Renderer = &Renderer{TemplateDir: "templates", Reload: debug, TemplateCache: make(map[string]*pongo2.Template)}

	e.Static("/static", "static")
	e.GET("/", index)
	e.GET("/about", about)
	e.GET("/login", getLogin)
	e.POST("/login", login)
	e.POST("/logout", logout)

	authedGroup := e.Group("/logbook", ensureSessionUser)
	authedGroup.GET("", getLogbook)
	authedGroup.GET("/add", getAddEntry)
	authedGroup.POST("/add", addEntry)
	authedGroup.GET("/:uuid", getEntry)
	authedGroup.POST("/:uuid", updateEntry)

	e.GET("/users/:username", notYetImplemented)
	e.GET("/users/:username/logbook", getUserLogbook)

	e.Logger.Fatal(e.Start(":8006"))
}
