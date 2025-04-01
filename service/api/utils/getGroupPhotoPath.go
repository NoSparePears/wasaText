package utils

import (
	"fmt"
)

func GetGroupPhotoPath(groupID int) string {
	return fmt.Sprintf("./storage/groups/%d.jpeg", groupID)
}
