package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var info map[int]*User
var rd *render.Render = render.New()
var lastID int

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func AddInformation(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		rd.JSON(w, http.StatusBadRequest, nil)
		return
	}

	// create
	lastID++
	user.ID = lastID
	info[user.ID] = user
	rd.JSON(w, http.StatusOK, info[user.ID])
}

func GetInformation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		rd.JSON(w, http.StatusBadRequest, nil)
		return
	}

	_, ok := info[id]
	if !ok {
		rd.JSON(w, http.StatusBadRequest, "No User")
		return
	}

	rd.JSON(w, http.StatusOK, info[id])
}

func DeleteInformation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		rd.JSON(w, http.StatusBadRequest, nil)
		return
	}

	_, ok := info[id]
	if !ok {
		rd.JSON(w, http.StatusBadRequest, "Can't Delete Because id already don't exists")
		return
	}

	delete(info, id)
	rd.JSON(w, http.StatusOK, "Delete Success")
}

func UpdateInformation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		rd.JSON(w, http.StatusBadRequest, nil)
		return
	}

	_, ok := info[id]

	if !ok {
		rd.JSON(w, http.StatusBadRequest, "Can't Update Because id don't exists")
		return
	}

	updatedUser := new(User)
	err = json.NewDecoder(r.Body).Decode(updatedUser)
	if err != nil {
		rd.JSON(w, http.StatusBadRequest, nil)
		return
	}

	updatedUser.ID = id
	info[id] = updatedUser
	rd.JSON(w, http.StatusOK, updatedUser)
}

func NewHandler() http.Handler {
	lastID = 0
	info = make(map[int]*User)

	router := mux.NewRouter()

	router.HandleFunc("/addinfo", AddInformation).Methods("POST")
	router.HandleFunc("/getinfo/{id:[0-9]+}", GetInformation).Methods("GET")
	router.HandleFunc("/deleteinfo/{id:[0-9]+}", DeleteInformation).Methods("DELETE")
	router.HandleFunc("/updateinfo/{id:[0-9]+}", UpdateInformation).Methods("PUT")

	router.PathPrefix("/css").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	n := negroni.Classic()
	n.UseHandler(router)

	return router
}
