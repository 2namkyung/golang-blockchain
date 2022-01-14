package myapp

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	response, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, response.StatusCode)

	data, _ := ioutil.ReadAll(response.Body)
	assert.Equal("Hello World", string(data))
}

func TestUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	response, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, response.StatusCode)

	data, _ := ioutil.ReadAll(response.Body)
	assert.Contains(string(data), "Get Userinfo")
}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	response, err := http.Get(ts.URL + "/users/100")
	assert.NoError(err)
	assert.Equal(http.StatusOK, response.StatusCode)

	data, _ := ioutil.ReadAll(response.Body)
	assert.Contains(string(data), "No User ID : 100")
}
