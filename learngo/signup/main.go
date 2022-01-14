package main

import (
	"learngo/signup/app"
	"net/http"
)

func main() {
	m, a := app.MakeHandler("./user.db")

	defer a.Close()

	http.ListenAndServe(":4000", m)
}
