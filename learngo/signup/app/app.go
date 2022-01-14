package app

import (
	"encoding/json"
	"fmt"
	"learngo/signup/model"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render = render.New()

type AppHandler struct {
	http.Handler
	db model.DBHandler
}

func (a *AppHandler) Close() {
	a.db.Close()
}

func (a *AppHandler) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userid := vars["userid"]
	fmt.Println(userid)
	user := a.db.GetUserInfo(userid)

	rd.JSON(w, http.StatusOK, user)
}

func (a *AppHandler) AddUserInfoHadler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rd.JSON(w, http.StatusBadRequest, "ERROR")
		return
	}

	userId := user.UserID
	name := user.Name
	password := user.Password
	birthDate := user.BirthDate
	gender := user.Gender
	phone := user.Phone
	location := user.Location

	NewUser := a.db.AddUserInfo(userId, name, password, birthDate, gender, phone, location)
	rd.JSON(w, http.StatusOK, NewUser)
}

func (a *AppHandler) RemoveUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userid := vars["userid"]
	check := a.db.RemoveUserInfo(userid)

	rd.JSON(w, http.StatusOK, check)
}

func (a *AppHandler) UpdateUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		rd.JSON(w, http.StatusBadRequest, "ERROR")
		return
	}

	// userId := user.UserID
	// name := user.Name
	// password := user.Password
	// birthDate := user.BirthDate
	// gender := user.Gender
	// phone := user.Phone
	// location := user.Location

	UpdatedUser := a.db.UpdateUserInfo(user.UserID)
	rd.JSON(w, http.StatusOK, UpdatedUser)
}

func MakeHandler(filepath string) (http.Handler, *AppHandler) {
	router := mux.NewRouter()
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(router)

	app := &AppHandler{Handler: n, db: model.NewDBHandler(filepath)}

	// Routing
	// User
	router.HandleFunc("/signup", app.AddUserInfoHadler).Methods("POST")
	router.HandleFunc("/user/{userid}", app.GetUserInfoHandler).Methods("GET")
	router.HandleFunc("/update", app.UpdateUserInfoHandler).Methods("UPDATE")
	router.HandleFunc("/delete/{userid}", app.RemoveUserInfoHandler).Methods("DELETE")

	// Item
	router.HandleFunc("/user_item/{userid}", app.GetItemInfoHandler).Methods("GET")
	router.HandleFunc("/addItem", app.AddItemInfoHandler).Methods("POST")
	router.HandleFunc("/updateItem", app.UpdateItemInfoHandler).Methods("UPDATE")
	router.HandleFunc("removeItem", app.RemoveItemInfoHandler).Methods("DELETE")

	// CORS
	credentials := handlers.AllowCredentials()
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Credentials"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	return handlers.CORS(headers, credentials, methods, origins)(app), app
}
