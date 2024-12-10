package api

import (
  "encoding/json"
  "net/http"
  "github.com/julienschmidt/httprouter"
  "wasaText/service/api/reqcontext"
)

func (rt *_router) CreateUser(u User) (User, error) {
  //create user in the db
  dbUser, err := rt.db.CreateUser(u.ConvertUserfromDB())
  if err != nil {
    return u,err
  }

  //convert user from db to api's user 
  err = u.ConvertUserfromDB(dbUser)
  if err != nil {
    return u, err
  }

  return u, nil 
}

func (rt *_router) doLogin (w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
  var user User
  
  //read request body
  err := json.NewDecoder(r.Body).Decode(&user)
  if err != nil {
    BadRequest(w, err, ctx, "Invalid username")
    return 
  }
  //check username 
  exists, err := rt.db.CheckIfExists(user.Username)
  if err != nil {
    InternalServerError(e, err, ctx)
    return 
  }
  if !exists {
    _, err := rt.CreateUser(user)
    if err != nil {
      InternalServerError(w, err, ctx)
      return 
    }
    w.WriteHeader(http.StatusCreated)
    
  } else {
    //it exists, so search it in the db 
    dbUser, err := rt.db.GetUserByName(user.Username)
    if err != nil {
      InternalServerError(w, err, ctx)
      return 
    }
    w.WriteHeader(http.StatusOK)
  }

  //Response 
  w.Header().Set("content-type", "application/json")
  if err := json.NewEncoder(w).Encode(user.ID); err != nil {
    InternalServerError(w, err, ctx)
    return 
  }
}
