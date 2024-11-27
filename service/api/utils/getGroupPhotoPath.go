package utils

import (
	"fmt"
)

func GetGroupPhotoPath(userID uint64, groupID uint64) string {
	return fmt.Sprintf("./storage/%d/groups/%d.jpeg", userID, groupID)
}
