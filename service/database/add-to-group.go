package database

var query_GROUP = `SELECT groupID from Group WHERE groupID = ? AND ownerID = ?`
var query_ADDMEMBER = `INSERT INTO GroupMember (groupID, memberID) VALUES (?, ?)`

func (db *appdbimpl) AddToGroup(user_id int, g_id int, added_user_id int) (Group, error) {
  var g Group
  _, err := db.c.Query(query_GROUP, g_id, user_id)
  if err != nil {
    return Group{}, err
  } else {
    _, err := db.c.Exec(query_ADDUSER, g_id, added_user_id)
    if err != nil {
      return Group{}, err
    }
  }
  return g, err
}
