package api

import(
	"regexp"
	"wasaText/service/api/utils"
	"wasaText/service/database"
)

//this struct represents the user object for API
type Conversation struct {
	ID int								`json:"id"`
	SendID int 						`json:"sendID"`
	RecID int							`json:"recID"`
	GroupID int						`json:"groupID"`
	LastMessageID int   	`json:"last_msgID"`
	DelBySend bool				`json:"del_by_send"`
	ConvoPropic64 string  `json:"photo"`
}

//converts the User strut to db representation
func (c *Conversation) ToDatabase() database.Conversation {
	return database.Conversation {
		ID: c.ID,
		SendID: c.SendID,
		RecID: c.RecID,
		GroupID: c.GroupID,
		LastMessageID: c.LastMessageID,
		DelBySend: c.DelBySend
		//profile picture not needed in db representation
	}
}

//populates the User struct with db data
func (c *Conversation) FromDatabase(dbConvo database.Conversation) error {
	c.ID = dbConvo.ID
	c.SendID = dbConvo.SendID
	c.RecID = dbConvo.RecID
	c.GroupID = dbConvo.GroupID
	c.LastMessageID = dbConvo.LastMessageID
	c.DelBySend = dbConvo.DelBySend
	//convert pfp path to Base64
	propic64, err := utils.ImageToBase64(utils.GetProfilePicPath(u.ID))
	if err != nil {
		return err 
	}

	c.ConvoPropic64 = propic64

	return nil 
}

//func isGroup????
