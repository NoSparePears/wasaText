package database

var query_GETCONVERSATION = `SELECT convoID FROM Conversation WHERE (sendID = ? AND recID = ?)`

func (db *appdbimpl) GetConversation(user_id int, id_rec int) (Conversation, error) {
  rows, err := db.c.Query(query_GETCONVERSATION, user_id, id_rec)
  
  if err != nil {
    return Conversation{}, err
  }
  defer func() { err = rows.Close() }()
  
  var convo Conversation

  if rows.Next() {
    if rows.Err() != nil {
      return Conversation{},err
    }

    err = rows.Scan(&convo.ID)
    if err != nil {
      return Conversation{}, err
    }
    
  }
  return convo, nil
}
