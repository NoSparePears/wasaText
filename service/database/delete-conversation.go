package database

var query_DELETECONVO = `DELETE FROM Conversation WHERE convoID = ? AND sendID = ?;`
var query_GETUSERCONVOS = `SELECT convoID FROM Conversation WHERE sendID = ?;`

func (db *appdbimpl) DeleteConversation(user_id int, convo_id int) ([]Conversation, error){
  result, err := db.c.Exec(query_DELETECONVO, convo_id, user_id)
  if err != nil {
    return nil, err
  }

  rowsAffected, err := result.RowsAffected()
  if err != nil {
    return nil, err
  }

  if rowsAffected == 0 {
    return nil, err
  }

  rows, err := db.c.Query(query_GETUSERCONVOS, user_id)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var convos []Conversation
  for rows.Next() {
    var convo Conversation
    if err := rows.Scan(&convo.ID); err != nil {
      return nil, err
    }
    convos = append(convos, convo)
  }

  if err := rows.Err(); err != nil {
    return nil, err
  }

  return convos, nil

  
}
