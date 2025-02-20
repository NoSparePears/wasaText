package database

// modifica come è scritto, non mi piace
var sql_USERTABLE = `CREATE TABLE IF NOT EXISTS User(
  userID INTEGER NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  PRIMARY KEY(userID)
  );`

// ok in teoria
var sql_GLOBALCONVOTABLE = `CREATE TABLE IF NOT EXISTS GlobalConversation (
  globalConvoID INTEGER PRIMARY KEY AUTOINCREMENT, 
  isGroup BOOLEAN DEFAULT FALSE
  );`

// ok in teoria
var sql_CONVOTABLE = `CREATE TABLE IF NOT EXISTS Conversation (
  userConvoID INTEGER PRIMARY KEY AUTOINCREMENT, 
  userID INTEGER NOT NULL,
  globalConvoID INTEGER NOT NULL,
  lastMsgId INTEGER,
  delByUser BOOLEAN DEFAULT FALSE,
  visible BOOLEAN DEFAULT TRUE,

  CONSTRAINT fk_user FOREIGN KEY (userID) REFERENCES User(userID),
  CONSTRAINT fk_global FOREIGN KEY (globalConvoID) REFERENCES GlobalConversation(globalConvoID),
  
  UNIQUE (userID, globalConvoID)
  );`

// ok in teoria
var sql_GROUPMEMBERTABLE = `CREATE TABLE IF NOT EXISTS GroupMember (
  groupID INTEGER NOT NULL,
  userID INTEGER NOT NULL.

  CONSTRAINT fk_group FOREIGN KEY (groupID) REFERENCES GlobalConversation(globalConvoID),
  CONSTRAINT fk_user FOREIGN KEY (userID) REFERENCES User(userID),

  UNIQUE (groupID, userID)
  );`

// ok in teoria
var sql_MSGTABLE = `CREATE TABLE IF NOT EXISTS Message (
  msgID INTEGER PRIMARY KEY AUTOINCREMENT,
  convoID NOT NULL,
  senderID NOT NULL,
  content TEXT NOT NULL,
  timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
  
  CONSTRAINT fk_convo FOREIGN KEY (convoID) REFERENCES GlobalConversation(globalConvoID),
  CONSTRAINT fk_user FOREIGN KEY (senderID) REFERENCES User(userID),
  
  UNIQUE (msgID, convoID)
);`

// rivedi
var sql_COMMTABLE = `CREATE TABLE IF NOT EXISTS Comment (
  commID INTEGER NOT NULL UNIQUE,
  msgID NOT NULL,
  sendID NOT NULL,
  emoji TEXT NOT NULL,
  PRIMARY KEY (commID, msgID),
  CONSTRAINT fk_comm
    FOREIGN KEY (msgID) REFERENCES Message(msgID)
      ON DELETE CASCADE,
    FOREIGN KEY (sendID) REFERENCES User(userID)    
  );`
