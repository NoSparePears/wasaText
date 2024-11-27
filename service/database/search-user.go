package database

var query_SEARCHUSER = `SELECT userID, username FROM User WHERE username = ?;`

func (db *appdbimpl) SearchUser(username string) (User, error) {
  var user User
  err := db.c.QueryRow(query_SEARCHUSER, username).Scan(&user.UserID, &user.Username)
  return user, err
}
