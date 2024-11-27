package database

import "time"

type User struct {
  ID uint64
  Username string 
}

type Conversation struct {
  ID uint64
  Messages []Message
}

type Message struct {
  ID uint64
  sent timestamp
  text string
}

type Group struct {
  ID uint64
  Participants []User
  Photo string 
}
