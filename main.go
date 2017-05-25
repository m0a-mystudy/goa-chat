//go:generate goagen bootstrap -d github.com/m0a-mystudy/goa-chat/design

package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	jwtgo "github.com/dgrijalva/jwt-go"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/security/jwt"
	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/controllers"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// ErrUnauthorized is the error returned for unauthorized requests.
	ErrUnauthorized = goa.NewErrorClass("unauthorized", 401)
)

func NewJWTMiddleware() goa.Middleware {
	keys, err := LoadJWTPublicKeys()
	if err != nil {
		panic(err)
	}
	return jwt.New(jwt.NewSimpleResolver(keys), ForceFail(), app.NewJWTSecurity())
}

func ForceFail() goa.Middleware {
	errValidationFailed := goa.NewErrorClass("validation_failed", 401)
	forceFail := func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
			if f, ok := req.URL.Query()["fail"]; ok {
				if f[0] == "true" {
					return errValidationFailed("forcing failure to illustrate Validation middleware")
				}
			}
			return h(ctx, rw, req)
		}
	}
	fm, _ := goa.NewMiddleware(forceFail)
	return fm
}

// LoadJWTPublicKeys loads PEM encoded RSA public keys used to validata and decrypt the JWT.
func LoadJWTPublicKeys() ([]jwt.Key, error) {
	keyFiles, err := filepath.Glob("./jwtkey/*.pub")
	if err != nil {
		return nil, err
	}
	keys := make([]jwt.Key, len(keyFiles))
	for i, keyFile := range keyFiles {
		pem, err := ioutil.ReadFile(keyFile)
		if err != nil {
			return nil, err
		}
		key, err := jwtgo.ParseRSAPublicKeyFromPEM([]byte(pem))
		if err != nil {
			return nil, fmt.Errorf("failed to load key %s: %s", keyFile, err)
		}
		keys[i] = key
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("couldn't load public keys for JWT security")
	}

	return keys, nil
}

func main() {
	// Create service
	service := goa.New("Chat API")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount security middlewares
	// app.UseBasicAuthMiddleware(service, NewBasicAuthMiddleware())
	app.UseJWTMiddleware(service, NewJWTMiddleware())

	user := os.Getenv("MYSQL_USER")
	if user == "" {
		service.LogError("startup", "err", fmt.Errorf("MYSQL_USER not found"))
		os.Exit(-1)

	}
	pass := os.Getenv("MYSQL_PASSWORD")
	if pass == "" {
		service.LogError("startup", "err", fmt.Errorf("MYSQL_PASSWORD not found"))
		os.Exit(-1)
	}

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@/goa_chat?parseTime=true", user, pass))
	if err != nil {
		service.LogError("startup", "err", err)
	}

	wsConns := controllers.NewConnections(service.Context)

	b, err := ioutil.ReadFile("./jwtkey/jwt.key")
	if err != nil {
		service.LogError("startup", "err", err)
		os.Exit(-1)

	}
	privKey, err := jwtgo.ParseRSAPrivateKeyFromPEM(b)
	if err != nil {
		service.LogError("startup", "err", fmt.Errorf("jwt: failed to load private key: %s", err))
		os.Exit(-1)
	}

	option := controllers.NewOption(db, wsConns, privKey)

	app.MountAccountController(service,
		controllers.NewAccountController(service, option))
	// Mount "message" controller
	app.MountMessageController(service,
		controllers.NewMessageController(service, option))
	// Mount "room" controller
	app.MountRoomController(service,
		controllers.NewRoomController(service, option))

	// Mount "serve" controller
	app.MountServeController(service,
		controllers.NewServeController(service))

	service.Mux.Handle("GET", "/login", func(w http.ResponseWriter, r *http.Request, _ url.Values) {
		controllers.MakeAuthHandlerFunc("")(w, r, option)
	})
	service.Mux.Handle("GET", "/oauth2callback", func(w http.ResponseWriter, r *http.Request, _ url.Values) {
		controllers.Oauth2callbackHandler(w, r, option)
	})
	// start service
	if err := service.ListenAndServe("oauth.local.com:8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}
