package database

var query_GETCONVERSATIONS = `SELECT convoID FROM Conversation WHERE sendID = ?;`

func (db *appdbimpl) GetConversations(user_id int) ([]Conversation, error) {
  //get conversations from database
  rows, err := db.c.Query(query_GETCONVERSATIONS, user_id)
  if err != nil {
    return nil, err
  }
  defer func() { err = rows.Close() }()

  var convos []Conversation

  for rows.Next() {
    if rows.Err() != nil {
      return nil, err
    }
    var convo Conversation

    err = rows.Scan(&convo.ID)
    if err != nil {
      return nil, err
    }
    convos = append(convos, convo)
  }

  return convos, err
}
