package model

func (s *sqliteHandler) GetUserInfo(userId string) User {

	var user User
	row := s.db.QueryRow("SELECT userId, name, password, birthDate, gender, phone, location FROM User WHERE userId=?", userId)

	row.Scan(&user.UserID, &user.Name, &user.Password, &user.BirthDate, &user.Gender, &user.Phone, &user.Location)
	return user
}

func (s *sqliteHandler) AddUserInfo(userId, name, passowrd, birthDate, gender, phone, location string) bool {
	statement, err := s.db.Prepare("INSERT INTO User(userId, name, password, birthDate, gender, phone, location) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}

	result, err := statement.Exec(userId, name, passowrd, birthDate, gender, phone, location)
	if err != nil {
		panic(err)
	}

	cnt, _ := result.RowsAffected()
	return cnt > 0
}

func (s *sqliteHandler) RemoveUserInfo(userId string) bool {

	statement, err := s.db.Prepare("DELETE FROM User WHERE userid=?")
	if err != nil {
		panic(err)
	}

	result, err := statement.Exec(userId)
	if err != nil {
		panic(err)
	}

	cnt, _ := result.RowsAffected()
	return cnt > 0
}

func (s *sqliteHandler) UpdateUserInfo(userId string) *User {
	user := &User{}
	row, err := s.db.Query("SELECT userId, name, password, birthDate, gender, phone, location FROM User=?", userId)
	if err != nil {
		panic(err)
	}

	defer row.Close()

	row.Scan(&user.UserID, &user.Name, &user.Password, &user.BirthDate, &user.Gender, &user.Phone, &user.Location)
	return user
}
