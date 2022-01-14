package app

import (
	"encoding/json"
	"learngo/signup/model"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *AppHandler) GetItemInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userid := vars["userid"]
	itemList := a.db.GetItemInfo(userid)
	rd.JSON(w, http.StatusOK, itemList)
}

func (a *AppHandler) AddItemInfoHandler(w http.ResponseWriter, r *http.Request) {
	var item model.Item

	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		rd.JSON(w, http.StatusBadRequest, "ERROR")
		return
	}

	name := item.Name
	price := item.Price
	market := item.Market
	createdAt := item.Created_at
	userId := item.UserId

	NewItem := a.db.AddItemInfo(name, price, market, createdAt, userId)
	rd.JSON(w, http.StatusOK, NewItem)
}

func (a *AppHandler) UpdateItemInfoHandler(w http.ResponseWriter, r *http.Request) {

}

func (a *AppHandler) RemoveItemInfoHandler(w http.ResponseWriter, r *http.Request) {

}
