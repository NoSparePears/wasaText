package database

func (db *appdbimpl) GetUserPhotoPath(userID int) (string, error) {
	var pfp64 string
	err := db.c.QueryRow("SELECT pfpPath FROM User WHERE userID = ?;", userID).Scan(&pfp64)
	if err != nil {
		return "", err
	}
	return pfp64, nil
}
