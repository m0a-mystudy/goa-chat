//go:generate goagen bootstrap -d github.com/m0a-mystudy/goa-chat/design

package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/controllers"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// ErrUnauthorized is the error returned for unauthorized requests.
	ErrUnauthorized = goa.NewErrorClass("unauthorized", 401)
)

// NewBasicAuthMiddleware creates a middleware that checks for the presence of a basic auth header
// and validates its content.
func NewBasicAuthMiddleware() goa.Middleware {
	return func(h goa.Handler) goa.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {

			goa.LogInfo(ctx, "header", fmt.Sprintf("<%s>", req.Header.Get("Authorization")))
			// Retrieve and log basic auth info
			user, pass, ok := req.BasicAuth()
			// A real app would do something more interesting here
			if !ok {
				goa.LogInfo(ctx, "failed basic auth")
				return ErrUnauthorized("missing auth")
			}
			if user != "abe" || pass != "pass" {
				return ErrUnauthorized("invalid auth")
			}

			// Proceed
			goa.LogInfo(ctx, "basic", "user", user, "pass", pass)
			return h(ctx, rw, req)
		}
	}
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
	app.UseBasicAuthMiddleware(service, NewBasicAuthMiddleware())

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
	// Mount "account" controller
	c3 := controllers.NewAccountController(service, db)
	app.MountAccountController(service, c3)

	// Mount "message" controller
	c := controllers.NewMessageController(service, db, wsConns)
	app.MountMessageController(service, c)
	// Mount "room" controller
	c2 := controllers.NewRoomController(service, db, wsConns)
	app.MountRoomController(service, c2)

	// Mount "serve" controller
	c4 := controllers.NewServeController(service)
	app.MountServeController(service, c4)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}
