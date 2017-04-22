package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// A User may log in to the site and create Entry records.
type User struct {
	ID          uint64    `db:"id"`
	UUID        string    `db:"uuid"`
	UserName    string    `db:"username"`
	Password    string    `db:"password"`
	DisplayName string    `db:"display_name"`
	Active      bool      `db:"active"`
	CreatedAt   time.Time `db:"created_at"`
	LastLoginAt time.Time `db:"last_login_at"`
}

// An Entry is a URL and/or Notes.
type Entry struct {
	ID        uint64         `db:"id"`
	UUID      string         `db:"uuid"`
	Title     string         `db:"title"`
	URL       string         `db:"url"`
	Notes     string         `db:"notes"`
	Private   bool           `db:"private"`
	UserID    uint64         `db:"user_id"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	Tags      pq.StringArray `db:"tags"`
}

var (
	db *sqlx.DB
)

// SetDB sets the database connection.
func SetDB(adb *sqlx.DB) {
	db = adb
}

var getUserLogbookSQL = `
	select   le.* 
	from     logbook_entry le 
				inner join logbook_user lu on (le.user_id = lu.id)
	where    lu.username = $1 
	and      le.private = false
	and      lu.active = true
	order by le.created_at desc
	limit $2
	offset $3
`

var getUserLogbookWithTagSQL = `
	select   le.* 
	from     logbook_entry le 
				inner join logbook_user lu on (le.user_id = lu.id)
	where    lu.username = $1
	and      $2 = any(le.tags) 
	and      le.private = false
	and      lu.active = true
	order by le.created_at desc
	limit $3
	offset $4
`

// GetUserLogbook returns an active user's public bookmarks.
func GetUserLogbook(username string, tag string, offset int, limit int) ([]*Entry, error) {
	var entries []*Entry

	if tag != "" {
		err := db.Select(&entries, getUserLogbookWithTagSQL, username, tag, limit, offset)
		return entries, err
	}

	err := db.Select(&entries, getUserLogbookSQL, username, limit, offset)
	return entries, err
}
