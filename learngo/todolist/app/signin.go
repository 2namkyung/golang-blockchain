package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserId struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback", // callback 받을 경로
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"},
	Endpoint:     google.Endpoint,
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	//googleOauthConfig.AuthCodeURL(State)
	// -> State is a token to protect the user from 'CSRF attacks'
	state := generateStateOauthCookie(w)

	// default : oauth2.AccessTypeOnline
	// if you want to use RefreshToken, use oauth2.AccessTypeOffline
	url := googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

	// 1. google url로 Redirect
	// 2. callback이 호출되었을 때,
	// Cookie의 oauthstate와 google에서 보내온 state와 동일하면 공격에 의한 로그인이 아니다
}

// 랜덤한 16바이트 값을 생성하고 Cookie의 Value로 만든다
func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour) // 만료시간 1일

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}

	http.SetCookie(w, cookie)

	return state
}

func googleoAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthstate.Value {
		errMsg := fmt.Sprintf("invalid google oauth state cookie : %s state ; %s\n", oauthstate.Value, r.FormValue("state"))
		log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store Id info into Session cookie
	// Unmarshal json
	var userInfo GoogleUserId
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Get a session
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Set some session values
	session.Values["id"] = userInfo.ID

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// fmt.Fprint(w, string(data))
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	// Context : 안전한 스레드간의 데이터 교환을 위한 저장소
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	fmt.Println("token : ", token)
	fmt.Println("RefreshToken : ", token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to Exchange %s", err.Error())
	}

	// Access Token , Refresh Token . . .
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to Get UserInfo %s", err.Error())
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
	//resp.Body의 ID를 User Key로 사용
}
