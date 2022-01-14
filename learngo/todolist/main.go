package main

import (
	"learngo/todolist/app"
	"net/http"
)

func main() {

	m := app.MakeHandler("./test.db")
	defer m.Close()

	http.ListenAndServe(":3000", m)

	/* Make UUID
	uuid := uuid.MakingUUID()
	fmt.Println(uuid)
	80fee3e3-49f2-405a-8e1e-c18a4a1bc71e
	*/
}
