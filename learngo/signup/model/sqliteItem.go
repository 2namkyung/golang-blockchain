package model

func (s *sqliteHandler) GetItemInfo(userId string) []Item {
	var itemList []Item

	rows, err := s.db.Query("SELECT name, price, market, created_at, updated_at FROM Item WHERE userId=?", userId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var item Item
		rows.Scan(&item.Name, &item.Price, &item.Market, &item.Created_at, &item.Updated_at)
		item.UserId = userId
		itemList = append(itemList, item)
	}

	return itemList
}

func (s *sqliteHandler) AddItemInfo(name, price, market, created_at, userId string) bool {
	statement, err := s.db.Prepare("INSERT INTO Item(name, price, market, created_at, userId) VALUES(?,?,?,?,?)")
	if err != nil {
		panic(err)
	}

	result, err := statement.Exec(name, price, market, created_at, userId)
	if err != nil {
		panic(err)
	}

	cnt, _ := result.RowsAffected()
	return cnt > 0
}

func (s *sqliteHandler) RemoveItemInfo(userId string) bool {

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

func (s *sqliteHandler) UpdateItemInfo(userId string) *Item {
	item := &Item{}
	row, err := s.db.Query("SELECT userId, name, password, birthDate, gender, phone, location FROM User=?", userId)
	if err != nil {
		panic(err)
	}

	defer row.Close()

	// row.Scan(&user.UserID, &user.Name, &user.Password, &user.BirthDate, &user.Gender, &user.Phone, &user.Location)
	return item
}
