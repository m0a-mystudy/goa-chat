//go:generate goagen bootstrap -d github.com/m0a-mystudy/goa-chat/design

package main

import (
	"database/sql"

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

	db, err := sql.Open("mysql", "root:sanset6@/goa_chat?parseTime=true")
	if err != nil {
		service.LogError("startup", "err", err)
	}
	// Mount "message" controller
	c := controllers.NewMessageController(service, db)
	app.MountMessageController(service, c)
	// Mount "room" controller
	c2 := controllers.NewRoomController(service, db)
	app.MountRoomController(service, c2)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}
