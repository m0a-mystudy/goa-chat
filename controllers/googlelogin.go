package controllers

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/m0a-mystudy/goa-chat/models"
	uuid "github.com/satori/go.uuid"
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

	mySigningKey = []byte("yasiayu2h8992js8")
)

// MakeAuthHandlerFunc return redirect
func MakeAuthHandlerFunc(redirectPath string) func(w http.ResponseWriter, r *http.Request, option *ControllerOptions) {
	return func(w http.ResponseWriter, r *http.Request, option *ControllerOptions) {
		if redirectPath == "" {
			redirectPath = "oauth2callback"
		}
		redierctURL := fmt.Sprintf("http://%s/%s", r.Host, redirectPath)
		conf.RedirectURL = redierctURL

		claims := &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(30) * time.Second).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		state, err := token.SignedString(mySigningKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		url := conf.AuthCodeURL(state)
		fmt.Println("url", url)
		http.Redirect(w, r, url, 302)
	}
}

func Oauth2callbackHandler(w http.ResponseWriter, r *http.Request, option *ControllerOptions) {

	state := r.FormValue("state")
	t, err := jwt.Parse(state, func(*jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if !t.Valid {
		http.Error(w, "state is invalid.", http.StatusUnauthorized)
		return
	}

	// 認証コードを取得します
	code := r.FormValue("code")
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

	// oauth2 clinet serviceを取得します
	service, err := v2.New(conf.Client(context, tok))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// token情報を取得します
	// ここにEmailやUser IDなどが入っています
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

	resp, err := http.Get(userInfo.Picture)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	defer resp.Body.Close()
	picture, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	googleUserID := tokenInfo.UserId
	account := &models.Account{}
	account, err = models.AccountByGoogleUserID(option.db, googleUserID)
	if err != nil {
		account = &models.Account{
			GoogleUserID: googleUserID,
			Image:        picture,
			Email:        userInfo.Email,
			Name:         userInfo.Name,
			Created:      time.Now(),
		}
		err = account.Insert(option.db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}

	token := jwt.New(jwt.SigningMethodRS512)
	in10m := time.Now().Add(time.Duration(10) * time.Minute).Unix()
	token.Claims = jwt.MapClaims{
		"iss":    "m0a",                 // who creates the token and signs it
		"exp":    in10m,                 // time when the token will expire (10 minutes from now)
		"jti":    uuid.NewV4().String(), // a unique identifier for the token
		"iat":    time.Now().Unix(),     // when the token was issued/created (now)
		"sub":    googleUserID,          // the subject/principal is whom the token is about
		"scopes": "api:access",          // token scope - not a standard claim
	}
	signedToken, err := token.SignedString(option.privateKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles("./html_template/save_token.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = tmpl.Execute(w, signedToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

}
