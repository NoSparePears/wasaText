package database

func (db *appdbimpl) GetGroupPhotoPath(groupID int) (string, error) {
	var pfp64 string
	err := db.c.QueryRow("SELECT photoPath FROM GlobalConversation WHERE globalConvoID = ?;", groupID).Scan(&pfp64)
	if err != nil {
		return "", err
	}
	return pfp64, nil
}
