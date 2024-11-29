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
	ChangeUsername(user_id int, name string) error
  
  //Search user by username
  SearchUser(username string) (User, error)

  //It returns an array of conversations 
  GetConversations(user_id int) ([]Conversation, error)

  //creates a new conversation in the database, given an user, it returns that conversation    
  CreateConversation(user_id int, id_rec int) (Conversation, error) 
	
  //It returns a specific conversation 
  GetConversation(user_id int, id_rec int) (Conversation, error)
    
  //It deletes a specific conversation
  DeleteConversation(user_id int, convo_id int) ([]Conversation, error)

  //It sends a Message
  SendMessage(user_id int, convo_id int, body string) (Message, error)
  
  //It deletes a Message
  DeleteMessage(user_id int, id int) (Conversation, error)

  //It comments a Message
  CommentMessage(user_id int, convo_id int, msg_id int, emoji string) (Message, error)

  //It forwards a massage to a conversation
  ForwardMessage(user_id int, msg_id int, convo_id int) (Conversation, error)

  //It deletes a comment 
  DeleteComment(user_id int, convo_id int, msg_id int, comm_id int) error
  
  //It creates a group
  CreateGroup(user_id int, name string, user_list []User) (Group, error)

  //It changes the group's name 
  SetGroupName(user_id int, g_id int, name string) (Group, error)

  //It changes the group's Photo
  SetGroupPhoto(user_id int, g_id int, photo string) (Group, error)

  //User leaves group
  LeaveGroup(user_id int, g_id int) ([]Conversation, error)
  
  //User adds another user to a group
  AddToGroup(user_id int, g_id int, added_user_id int) (Group, error)


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
