//go:generate goagen bootstrap -d github.com/m0a-mystudy/goa-chat/design

package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/controllers"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Create service
	service := goa.New("Chat API")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

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
	// Mount "message" controller
	c := controllers.NewMessageController(service, db, wsConns)
	app.MountMessageController(service, c)
	// Mount "room" controller
	c2 := controllers.NewRoomController(service, db, wsConns)
	app.MountRoomController(service, c2)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}
