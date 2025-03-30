package database

func (db *appdbimpl) UpdateUserPhotoPath(userID int, photoPath string) error {
	_, err := db.c.Exec("UPDATE User SET pfpPath = ? WHERE userID = ?", photoPath, userID)
	return err
}
