package api

func (rt *_router) CreateUser(u User) (User, error) {
	//return user object with his identificator
	dbUser, err := rt.db.CreateUser(u.ToDatabase())
	if err != nil {
		return u, err 
	}

	//checks if the user has been succesfully added to the db
	err = u.FromDatabase(dbUser)
	if err != nil {
		return u, err 
	}

	return u, nil 

}
