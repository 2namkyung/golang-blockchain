package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var userMap map[int]*User
var lastID int

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get Userinfo by /users/{id}")
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	user, ok := userMap[id]

	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID : ", id)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.MarshalIndent(user, "", "")
	fmt.Fprint(w, string(data))
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	//Create User
	lastID++
	user.ID = lastID
	user.CreatedAt = time.Now()
	userMap[user.ID] = user

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.MarshalIndent(user, "", "")
	fmt.Fprint(w, string(data))
}

func GetUserInfoAll(w http.ResponseWriter, r *http.Request) {
	for index, value := range userMap {
		data, _ := json.MarshalIndent(value, "", "")
		fmt.Fprint(w, "number : ", index, "\n", string(data), "\n\n")
	}
}

// Delete UserInfo in userMap
func DeleteUserInfo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	_, ok := userMap[id]

	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID : ", id, "So Can't delete that UserInfo")
		return
	}

	delete(userMap, id)
}

// NewHttpHandler make a new myapp handler
func NewHttpHandler() http.Handler {
	userMap = make(map[int]*User)
	lastID = 0

	/* gorilla.mux */
	mux := mux.NewRouter()

	mux.HandleFunc("/", IndexHandler)
	mux.Handle("/", http.FileServer(http.Dir("public")))

	mux.HandleFunc("/users", UsersHandler).Methods("GET")
	mux.HandleFunc("/users", CreateUserHandler).Methods("POST")
	mux.HandleFunc("/users/{id:[0-9]+}", GetUserInfoHandler)
	mux.HandleFunc("/usersAll", GetUserInfoAll)
	mux.HandleFunc("/userDelete/{id:[0-9]+}", DeleteUserInfo)

	/* goriila.pat */
	//mux := pat.New()
	// mux.Get("/users", UsersHandler)
	// mux.Post("/users", CreateUserHandler)

	// negroni : http middleware
	n := negroni.Classic()
	n.UseHandler(mux)
	return mux

}
