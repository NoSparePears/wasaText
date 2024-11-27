/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrUserDoesNotExist = errors.New("user does not exists")

// AppDatabase is the high level interface for the DB
type AppDatabase interface {

  //SetName allows to set/change username 
	SetName(name string) error
  
  //It returns an array of conversations 
  ListConversations() ([]Conversation, error)

  //creates a new conversation in the database, given an user, it returns that conversation    
  CreateConversation(id_rec uint64) (Conversation, error) 
	
  //It returns a specific conversation 
  GetConversation(id_rec uint64) (Conversation, error)
    
  //It deletes a specific conversation
  DeleteConversation(id uint64) ([]Conversation, error)

  //It sends a Message
  SendMessage(id uint64, body string) (Message, error)
  
  //It deletes a Message
  DeleteMessage(id uint64) (Conversation, error)

  //It comments a Message
  CommentMessage(msg_id uint64, emoji string) (Message, error)

  //It forwards a massage to a conversation
  ForwardMessage(msg_id uint64, convo_id uint64) (Conversation, error)

  //It deletes a comment 
  DeleteComment(comm_id uint64) error
  
  //It creates a group
  CreateGroup(name string, user_list []User) (Group, error)

  //It changes the group's name 
  SetGroupName(g_id intui64, name string) (Group, error)

  //It changes the group's Photo
  SetGroupPhoto(g_id intui64, photo string) (Group, error)

  //User leaves group
  LeaveGroup(g_id intui64) ([]Conversation, error)
  
  //User adds another user to a group
  AddToGroup(g_id intui64, user_id intui64) (Group, error)


  Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
