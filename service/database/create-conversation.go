package database

import (
  "database/sql"
  "errors"
  "wasaText/service/api/utils"
)

var query_CREATECONVO = `INSERT INTO Conversation (convoID, sendID, recID) VALUES (?, ?, ?)`
var query_CREATEGROUP = `INSERT INTO Conversation (convoID, sendID, groupID) VALUES (?, ?, ?)`

var query_LASTCONVOID = `SELECT MAX(convoID) FROM Conversation WHERE userID = ?`

func GetLastConvoId (db *appdbimpl, user_id int) (int, error) {
  var _maxID = sql.NullInt64{Int64: 0, Valid: false}
  row, err := db.c.Query(query_LASTCONVOID, user_id)
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

func (db *appdbimpl) CreateConversation(user_id int, id_rec int, id_group int) (Conversation, error)  {
  var convo Conversation
  
  maxID, err := GetLastConvoId (db, user_id)
  if err != nil {
    return convo, err
  }
  convo.ID = maxID + 1
  
  if id_group != 0 {
    convo.GroupID = id_group
    //aggiungi membri
  } else {
    convo.RecID = id_rec
    //inserisco nel db
    _, err := db.c.Exec(query_CREATECONVO, convo.ID, user_id, id_rec)
    if err != nil {
      return convo, err
    }
  }

  return convo, nil
}
