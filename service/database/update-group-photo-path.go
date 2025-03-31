package database

func (db *appdbimpl) UpdateGroupPhotoPath(groupID int, photoPath string) error {
	_, err := db.c.Exec("UPDATE GlobalConversation SET photoPath = ? WHERE globalConvoID = ?", photoPath, groupID)
	return err
}
