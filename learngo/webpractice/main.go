package main

import (
	"learngo/webpractice/api"
	"net/http"
)

func main() {

	http.ListenAndServe(":1234", api.NewHandler())
}
