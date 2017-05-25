package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	v2 "google.golang.org/api/oauth2/v2"
)

var (
	conf = oauth2.Config{
		ClientID:     os.Getenv("OPENID_GOOGLE_CLIENT"),
		ClientSecret: os.Getenv("OPENID_GOOGLE_SECRET"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL: "https://www.googleapis.com/oauth2/v4/token",
		},
	}

	state = "testSampleCodeStae"
)

func getRedirectURL(r *http.Request) string {
	return fmt.Sprintf("http://%s/oauth2callback", r.Host)
}

// MakeAuthHandlerFunc return redirect
func MakeAuthHandlerFunc(redirectPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if redirectPath == "" {
			redirectPath = "oauth2callback"
		}
		redierctURL := fmt.Sprintf("http://%s/%s", r.Host, redirectPath)
		getRedirectURL(r)
		conf.RedirectURL = redierctURL
		url := conf.AuthCodeURL(state)
		fmt.Println("url", url)
		http.Redirect(w, r, url, 302)
	}
}

func Oauth2callbackHandler(w http.ResponseWriter, r *http.Request) {

	if r.FormValue("state") != state {
		http.Error(w, "state is invalid.", http.StatusUnauthorized)
		return
	}

	fmt.Println(r.FormValue("state"))

	// 認証コードを取得します
	code := r.FormValue("code")
	// appengineのcontextを取得します
	// context := appengine.NewContext(r)
	context := context.Background()
	// 認証コードからtokenを取得します
	tok, err := conf.Exchange(context, code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// tokenが正しいことを確認します
	if tok.Valid() == false {
		http.Error(w, "token is invalid.", http.StatusUnauthorized)
		return
	}

	//tok.RefreshToken

	// oauth2 clinet serviceを取得します
	// 特にuserの情報が必要ない場合はスルーです
	service, err := v2.New(conf.Client(context, tok))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// token情報を取得します
	// ここにEmailやUser IDなどが入っています
	// 特にuserの情報が必要ない場合はスルーです
	tokenInfo, err := service.Tokeninfo().AccessToken(tok.AccessToken).Context(context).Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userInfo, err := service.Userinfo.Get().Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	fmt.Println(tokenInfo.UserId)
	fmt.Println(userInfo.Picture)
	fmt.Println(userInfo.Email)
}
