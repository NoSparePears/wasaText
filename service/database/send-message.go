package database

import (
  "database/sql"
  "errors"
)

var query_INSERTMSG = `INSERT INTO Message (msgID, convoID, sendID, content) VALUES (?, ?, ?, ?)`
var query_ID = `SELECT MAX(msgID) FROM Message WHERE convoID = ?`

func GetLastMsgId (db *appdbimpl, convo_id int) (int, error) {
  var _maxID = sql.NullInt64{Int64: 0, Valid: false}
  row, err := db.c.Query(query_ID, convo_id)
  if err != nil {
    return 0, err
  }
  var maxID int
  for row.Next() {
    if row.Err() != nil {
      return 0, err
    }
    err = row.Scan(&_maxID)
    if err != nil && !errors.Is(err, sql.ErrNoRows) {
      return 0, err
    }
    if !_maxID.Valid {
      maxID = 0
    } else {
      maxID = int(_maxID.Int64)
    }
  }
  return maxID, nil
}

func (db *appdbimpl) SendMessage(user_id int, convo_id int, body string) (Message, error){
  var msg Message
  
  maxID, err := GetLastMsgId (db, convo_id)
  if err != nil {
    return msg, err
  }

  msg.ID = maxID + 1

  _, err = db.c.Exec(query_INSERTMSG, msg.ID, convo_id, user_id, body)
  if err != nil {
    return msg, err
  }
  
  return msg, nil
}
