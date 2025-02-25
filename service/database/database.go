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
	"context"
	"database/sql"
	"errors"
	"fmt"
	"wasaText/service/structs"
)

var ErrUserDoesNotExist = errors.New("user does not exists")

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(username string) (structs.User, error)

	//ChangeUsername allows to set/change username
	ChangeUsername(user_id int, name string) error

	//Search user by username

	SearchUser(username string) (structs.User, error)

	//It returns an array of conversations
	GetConversations(user_id int) ([]structs.Conversation, error)

	//creates a new conversation in the database, given an user, it returns that conversation
	CreateConversation(user_id int, rec_id int) (structs.Conversation, error)

	//It returns a specific conversation
	GetConversation(user_id int, id_rec int) (structs.Conversation, error)

	//It deletes a specific conversation
	DeleteConversation(user_id int, convo_id int) ([]structs.Conversation, error)

	//It sends a Message
	InsertMessage(msg structs.Message, recID int) (structs.Message, error)

	//It deletes a Message
	DeleteMessage(msgID int, convoID int, senderID int) (structs.Conversation, error)

	//It comments a Message
	InsertComment(comment structs.Comment) (structs.Comment, error)

	//It forwards a massage to a conversation
	ComposeMsgToForward(msgID int, convoID int, customContent string) (structs.Message, error)

	//It deletes a comment
	DeleteComment(user_id int, msg_id int, comm_id int) error

	//It creates a group
	CreateGroup(groupName string, userID int) (structs.Group, error)

	//It adds another user to a group
	AddToGroup(user_id int, g_id int) error

	//It returns a list for all the members in a group
	GetGroupMembers(groupID int) ([]structs.User, error)

	//It changes the group's name
	SetGroupName(groupID int, name string) error

	//User leaves group
	LeaveGroup(userID int, groupID int) error

	//It adds checkmark for a sent message, and its timestamp
	AddSentCheck(msgID int) error

	Ping() error
}

type appdbimpl struct {
	c   *sql.DB
	ctx context.Context
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	//check if the database is nil (required)
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableCount uint8
	err := db.QueryRow(`SELECT COUNT(name) FROM sqlite_master WHERE type='table';`).Scan(&tableCount)
	if errors.Is(err, sql.ErrNoRows) {
		if err != nil {
			return nil, fmt.Errorf("error checking if database is empty: %w", err)
		}
	}

	// The tables are six in total, so if the count is less than 6, we need to create them.
	if tableCount != 7 {

		// ---CREATE USER TABLE----//
		_, err = db.Exec(sql_USERTABLE)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure user: %w", err)
		}

		// ---CREATE GLOBAL CONVO TABLE----//
		_, err = db.Exec(sql_GLOBALCONVOTABLE)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure global conversation: %w", err)
		}

		// ---CREATE CONVO TABLE----//
		_, err = db.Exec(sql_CONVOTABLE)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure conversation: %w", err)
		}

		// ---CREATE GROUP MEMBER TABLE----//
		_, err = db.Exec(sql_GROUPMEMBERTABLE)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure group member: %w", err)
		}

		// ---CREATE MESSAGE TABLE----//
		_, err = db.Exec(sql_MSGTABLE)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure message: %w", err)
		}

		// ---CREATE COMMENT TABLE----//
		_, err = db.Exec(sql_COMMTABLE)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure comment: %w", err)
		}

		// ---CREATE CHECKMARKS TABLE----//
		_, err = db.Exec(sql_CHECKMARKSTABLE)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure checkmarks: %w", err)
		}

	}
	return &appdbimpl{
		c:   db,
		ctx: context.Background(),
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
