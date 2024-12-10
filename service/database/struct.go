package database

import "time"

type User struct {
	ID int						`json:"id"`
	Username string 	`json:"username"`
}

type Conversation struct { 
	ID int							`json:"id"`
	SendID int 					`json:"sendID"`
	RecID int						`json:"recID"`
	GroupID int					`json:"groupID"`
	LastMessageID int		`json:"last_msgID"`
	DelBySend bool			`json:"del_by_send"`
}

type Message struct {
	ID int							`json:"id"`
	ConvoID int					`json:"convoID"`
	Timestamp time.Time	`json:"timestamp"`
	Body string					`json:"body"`
}

type Group struct {
	ID int 					`json:"id"`
	Members []User	`json:"members"`
	Name string			`json:"name"`
}
