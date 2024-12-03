package database

import "time"

type User struct {
  ID int
  Username string 
}

type Conversation struct { 
  ID int
  SendID int 
  RecID int
  GroupID int
  Messages []Message
  DelBySend bool
}

type Message struct {
  ID int
  ConvoID int
  Timestamp time.Time 
  text string
}

type Group struct {
  ID int 
  Members []User
  Name string
  Photo string 
}
