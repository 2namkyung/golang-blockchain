package app

import (
	"fmt"
	"io/ioutil"
	"learngo/todolist/model"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/skip2/go-qrcode"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var rd *render.Render = render.New()
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type AppHandler struct {
	http.Handler //Embeded
	db           model.DBHandler
}

type Success struct {
	Success bool `json:"success"`
}

// testing을 위한 로그인 여부 확인 함수 변경
// -> function pointer를 값으로 가지는 getSessionID 변수로 변경
var getSessionID = func(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	val := session.Values["id"]
	if val == nil {
		return ""
	}

	return val.(string)
}

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	list := a.db.GetTodos(sessionId)
	rd.JSON(w, http.StatusCreated, list)
}

func (a *AppHandler) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	name := r.FormValue("name")
	todo := a.db.AddTodo(name, sessionId)

	rd.JSON(w, http.StatusCreated, todo)
}

func (a *AppHandler) removeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	ok := a.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}
}

func (a *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"

	ok := a.db.CompleteTodo(id, complete)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusOK, Success{false})
	}

}

func (a *AppHandler) Close() {
	a.db.Close()
}

func CheckSigin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// if request URL is /signin.html, then next()
	if strings.Contains(r.URL.Path, "/signin") ||
		strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		return
	}

	// if user already signed in
	sessionID := getSessionID(r)
	if sessionID != "" {
		next(w, r)
		return
	}

	// if not user signed in , redirect signin.html
	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)

}

func QrCode(w http.ResponseWriter, r *http.Request) {

	err := qrcode.WriteFile("http://localhost:3000/todo.html", qrcode.Medium, 256, "test.png")
	if err != nil {
		fmt.Println(err)
	}

	fileBytes, err := ioutil.ReadFile("test.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	rd.JSON(w, http.StatusOK, fileBytes)
}

func MakeHandler(filepath string) *AppHandler {
	// model.TodoMap = make(map[int]*model.Todo)

	router := mux.NewRouter()
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.HandlerFunc(CheckSigin), negroni.NewStatic(http.Dir("public"))) // 3가지의 Decorator : NewRecovery(), NewLogger(), NewStatic(http.Dir("public"))                                                                                                     //add middleware
	n.UseHandler(router)

	fmt.Println(http.Dir("public"))

	app := &AppHandler{Handler: n, db: model.NewDBHandler(filepath)}

	router.HandleFunc("/", app.indexHandler)
	router.HandleFunc("/todos", app.getTodoListHandler).Methods("GET")
	router.HandleFunc("/todos", app.addTodoHandler).Methods("POST")
	router.HandleFunc("/todos/{id:[0-9]+}", app.removeTodoHandler).Methods("DELETE")
	router.HandleFunc("/complete-todo/{id:[0-9]+}", app.completeTodoHandler).Methods("GET")
	router.HandleFunc("/qrcode", QrCode).Methods("GET")

	// oAuth Google
	router.HandleFunc("/auth/google/login", googleLoginHandler)
	router.HandleFunc("/auth/google/callback", googleoAuthCallback)
	return app
}
