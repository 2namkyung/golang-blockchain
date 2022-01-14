package banking

import "errors"

// Account struct - 대문자 Public , 소문자 Private
type Account struct {
	owner   string
	balance int
}

var errNoMoney = errors.New("can't withdraw you are pool")

//NewAccount : Create Account
func NewAccount(owner string) *Account {
	account := Account{owner: owner, balance: 0}
	return &account
}

//Deposit Money
//Reciever는 Value Copy이므로 (a *Account) 사용
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

//Balance : Get Balance
func (a Account) Balance() int {
	return a.balance
}

//Withdraw Money
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}

//ChangeOwner : NewOwner
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

//Owner of Account
func (a Account) Owner() string {
	return a.owner
}

func (a Account) String() string {
	return "Whatever you want"
}
