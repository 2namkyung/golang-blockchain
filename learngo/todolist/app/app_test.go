package app

import (
	"encoding/json"
	"fmt"
	"learngo/todolist/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTodos(t *testing.T) {
	os.Remove("./test.db")

	assert := assert.New(t)

	getSessionID = func(r *http.Request) string {
		return "testsessionId"
	}

	//AppHandler는 http.Handler를 Embeded하고 있기 때문에 다음과 같이 사용할 수 있다.
	ah := MakeHandler("./test.db")
	defer ah.Close()

	ts := httptest.NewServer(ah)
	defer ts.Close()

	// Data add시 JSON을 보낸 것이 아니라 FormValue을 전송하였기 때문에 PostForm사용
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	/* Add & Get Test */
	var todo model.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")

	id1 := todo.Id
	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name": {"Test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo2")

	id2 := todo.Id

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	todos := []*model.Todo{}
	json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.Id == id1 {
			assert.Equal("Test todo", t.Name)
		} else if t.Id == id2 {
			assert.Equal("Test todo2", t.Name)
		} else {
			assert.Error(fmt.Errorf("testId should be id1 or id2"))
		}
	}

	/* Complete Test */
	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	todos = []*model.Todo{}
	json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	for _, t := range todos {
		if t.Id == id1 {
			assert.True(t.Completed)
		}
	}

	/* Delete Test */
	// Golang http pkg는 GET, POST는 지원하지만 DELETE는 지원하지 않는다
	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(id1), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	todos = []*model.Todo{}
	json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)
	assert.Equal(todo.Id, id2)
}
