package dictionary

import "errors"

//Dictionary type -> Using for map as Method Reciever
type Dictionary map[string]string // Dictionary = map[string]string 's alias
//map[key] -> return value, ok(boolean)

var errNotFound = errors.New("not Found")
var errWordExists = errors.New("that word already Exists")

//Search for a word
func (d Dictionary) Search(word string) (string, error) {
	value, exists := d[word]
	if exists {
		return value, nil
	}
	return "", errNotFound
}

//Add word
func (d Dictionary) Add(word, def string) error {
	_, err := d.Search(word)
	switch err {
	case errNotFound:
		d[word] = def
	case nil:
		return errWordExists
	}
	return nil
}

/*
	main()
	-> a := Dictionary{"Key":"Value"}
	value, err := a.Search("Key")
	if err!=nil -> print(err)
	else -> print(value)
*/
