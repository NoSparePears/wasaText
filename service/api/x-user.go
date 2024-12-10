package api

import(
	"regexp"
	"wasaText/service/api/utils"
	"wasaText/service/database"
)

//this struct represents the user object for API
type User struct {
	ID int							`json:"id"`
	Username string			`json:"username"`
	UserPropic64 string `json:"photo"` //Base64 encode pfp
}

//converts the User strut to db representation
func (u *User) ToDatabase() database.User {
	return database.User{
		ID: u.ID,
		Username: u.Username
		//profile picture not needed in db representation
	}
}

//populates the User struct with db data
func (u *User) FromDatabase(dbUser database.User) error {
	u.ID = dbUser.ID
	u.Username = dbUser.Username

	//convert pfp path to BAse64
	propic64, err := utils.ImageToBase64(utils.GetProfilePicPath(u.ID))
	if err != nil {
		return err 
	}

	u.UserPropic64 = propic64

	return nil 
}

//validates username according to set rules
func (u *User) IsValid() bool {
	username := u.Username
	validUser := regexp.MustCompile(`[a-z][a-z0-9]{2,13}$`)
	return validUser.MatchString(username)

}
