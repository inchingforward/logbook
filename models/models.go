package models

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// A User may log in to the site and create Entry records.
type User struct {
	ID          uint64      `db:"id"`
	UUID        string      `db:"uuid"`
	UserName    string      `db:"username"`
	Password    string      `db:"password"`
	DisplayName string      `db:"display_name"`
	Active      bool        `db:"active"`
	CreatedAt   time.Time   `db:"created_at"`
	LastLoginAt pq.NullTime `db:"last_login_at"`
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

const getUserLogbookSQL = `
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

const getUserLogbookWithTagSQL = `
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

const getLogbookSQL = `
	select   le.* 
	from     logbook_entry le 
			 inner join logbook_user lu on (le.user_id = lu.id)
	where    lu.id = $1
	and      lu.active = true
	order by le.created_at desc
	limit $2
	offset $3
`

const getLogbookWithTagSQL = `
	select   le.* 
	from     logbook_entry le 
			 inner join logbook_user lu on (le.user_id = lu.id)
	where    lu.id = $1
	and      $2 = any (le.tags)
	and      lu.active = true
	order by le.created_at desc
	limit $3
	offset $4
`

const getLogbookEntry = `
	select *
	from   logbook_entry
	where  user_id = $1
	and    uuid = $2
`

// GetUserPublicLogbook returns an active user's public bookmarks.
func GetUserPublicLogbook(username string, tag string, offset int, limit int) ([]*Entry, error) {
	var entries []*Entry

	if tag != "" {
		err := db.Select(&entries, getUserLogbookWithTagSQL, username, tag, limit, offset)
		return entries, err
	}

	err := db.Select(&entries, getUserLogbookSQL, username, limit, offset)
	return entries, err
}

// GetLogbook returns a logged-in user's logbook.
func GetLogbook(userID uint64, tag string, offset int, limit int) ([]*Entry, error) {
	var entries []*Entry

	if tag != "" {
		err := db.Select(&entries, getLogbookWithTagSQL, userID, tag, limit, offset)
		return entries, err
	}

	err := db.Select(&entries, getLogbookSQL, userID, limit, offset)

	return entries, err
}

// Login returns a User if the username and password match an existing user.
func Login(username string, password string) (User, error) {
	var user User

	err := db.Get(&user, "select * from logbook_user where username = $1", username)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, err
	}

	return user, err
}

// GetLogbookEntry returns an unique Entry by user id and entry uuid.
func GetLogbookEntry(userID uint64, entryUUID string) (Entry, error) {
	var entry Entry

	err := db.Get(&entry, getLogbookEntry, userID, entryUUID)

	return entry, err
}
