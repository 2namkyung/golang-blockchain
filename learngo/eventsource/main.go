package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/antage/eventsource"
	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
)

type Message struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

var msgCh chan Message

func PostMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	name := r.FormValue("name")

	SendMessage(name, msg)
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	SendMessage("", fmt.Sprintf("add user :  %s", username))
}

func SendMessage(name, msg string) {
	// send msg to every client
	msgCh <- Message{name, msg}
}

func ProcessMsgCh(es eventsource.EventSource) {
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond()))
	}
}

func LeftUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	SendMessage("", fmt.Sprintf("left user : %s", username))
}

func main() {

	msgCh = make(chan Message)

	es := eventsource.New(nil, nil)
	defer es.Close()

	go ProcessMsgCh(es)

	mux := pat.New()
	mux.Post("/messages", PostMessageHandler)
	mux.Handle("/stream", es)
	mux.Post("/users", AddUserHandler)
	mux.Delete("/users", LeftUserHandler)

	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
