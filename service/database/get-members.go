package database

var query_GETMEMBERS = `SELECT memberID FROM GroupMember WHERE groupID = ?`

func (db *appdbimpl) GetGroupMembers(group_id int) ([]User, error){
  var members []User

  rows, err := db.c.Query(query_GETMEMBERS, group_id)
  if err != nil {
    return 0, err
  }
  
  for rows.Next(){
    if rows.Err() != nil {
      return nil, err
    }
    var member User
    err = rows.Scan(&User.ID)
    if err != nil {
      return nil, err
    }
    convos = append(members, member)
  }
  return members, err
}
