package model

type User struct {
	UserID    string `json:"user_id"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	BirthDate string `json:"birth_date"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Location  string `json:"location"`
}

type Item struct {
	Name       string `json:"name"`
	Price      string `json:"price"`
	Market     string `json:"market"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	UserId     string `json:"user_id"`
}

type DBHandler interface {

	// User
	GetUserInfo(userId string) User
	AddUserInfo(userId, name, password, birthDate, gender, phone, location string) bool
	RemoveUserInfo(userId string) bool
	UpdateUserInfo(userId string) *User

	// Item
	GetItemInfo(userId string) []Item
	AddItemInfo(name, price, market, created_at, userId string) bool
	RemoveItemInfo(userId string) bool
	UpdateItemInfo(userId string) *Item
	Close()
}

func NewDBHandler(filepath string) DBHandler {
	return newSqliteHandler(filepath)
}
