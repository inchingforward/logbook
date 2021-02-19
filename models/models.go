package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
	ID        uint64         `db:"id" form:"id"`
	UUID      string         `db:"uuid" form:"uuid"`
	Title     string         `db:"title" form:"title"`
	URL       string         `db:"url" form:"url"`
	Notes     string         `db:"notes" form:"notes"`
	Private   bool           `db:"private" form:"private"`
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

// Validate validates a Logbook Entry, returning an error if it is not valid.
func (entry *Entry) Validate() error {
	if entry.Title == "" {
		return errors.New("title is required")
	}

	return nil
}

// NamedInsert executes the query insert statement and returns the
// generated sequence id.
func NamedInsert(query string, arg interface{}) (uint64, error) {
	rows, err := db.NamedQuery(query, arg)
	if err != nil {
		return 0, err
	}

	if !rows.Next() {
		return 0, rows.Err()
	}

	var id uint64
	err = rows.Scan(&id)
	rows.Close()
	if err != nil {
		return 0, err
	}

	return id, nil
}

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

// InsertEntry saves a new Entry into the database.
func InsertEntry(userID uint64, entry *Entry) error {
	entry.UUID = uuid.New().String()
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	entry.UserID = userID

	id, err := NamedInsert(`
		insert into logbook_entry 
		values (default, :uuid, :title, :url, :notes, :private, :user_id, :created_at, :updated_at, :tags) 
		returning id
	`, entry)

	if err != nil {
		return err
	}

	entry.ID = id

	return nil
}

// UpdateEntry updates a logbook entry in the database.
func UpdateEntry(entry *Entry) error {
	entry.UpdatedAt = time.Now()

	_, err := db.NamedQuery(`
		update logbook_entry 
		set    title = :title, 
		       url = :url, 
			   notes = :notes, 
			   private = :private,
			   updated_at = :updated_at,
			   tags = :tags
		where  id = :id
	`, entry)

	return err
}
