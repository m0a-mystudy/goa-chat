package controllers

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
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

	state = "testSampleCodeState"
)

func getRedirectURL(r *http.Request) string {
	return fmt.Sprintf("http://%s/oauth2callback", r.Host)
}

// MakeAuthHandlerFunc return redirect
func MakeAuthHandlerFunc(redirectPath string) func(w http.ResponseWriter, r *http.Request, option *ControllerOptions) {
	return func(w http.ResponseWriter, r *http.Request, option *ControllerOptions) {
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

func Oauth2callbackHandler(w http.ResponseWriter, r *http.Request, option *ControllerOptions) {

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

	//105212080826678848901
	fmt.Println(tokenInfo.UserId)
	fmt.Println(userInfo.Picture)

	resp, _ := http.Get(userInfo.Picture)
	defer resp.Body.Close()
	picture, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(len(picture))

	fmt.Println(userInfo.Name)
	fmt.Println(userInfo.Email)

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

	fmt.Printf("account = %v", account)

	token := jwtgo.New(jwtgo.SigningMethodRS512)
	in10m := time.Now().Add(time.Duration(10) * time.Minute).Unix()
	token.Claims = jwtgo.MapClaims{
		"iss":      "Issuer",              // who creates the token and signs it
		"aud":      "Audience",            // to whom the token is intended to be sent
		"exp":      in10m,                 // time when the token will expire (10 minutes from now)
		"jti":      uuid.NewV4().String(), // a unique identifier for the token
		"iat":      time.Now().Unix(),     // when the token was issued/created (now)
		"nbf":      2,                     // time before which the token is not yet valid (2 minutes ago)
		"sub":      "subject",             // the subject/principal is whom the token is about
		"scopes":   "api:access",          // token scope - not a standard claim
		"googleID": googleUserID,
	}
	signedToken, err := token.SignedString(option.privateKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set auth header for client retrieval
	resp.Header.Set("Authorization", "Bearer "+signedToken)

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
