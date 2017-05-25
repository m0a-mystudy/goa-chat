package controllers

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware/security/jwt"

	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/models"
)

// ToRoomMedia convert tool
func ToRoomMedia(room *models.Room) *app.Room {
	ret := app.Room{
		ID:          &room.ID,
		Description: room.Description,
		Name:        room.Name,
		Created:     &room.Created,
	}
	return &ret
}

// RoomController implements the room resource.
type RoomController struct {
	*goa.Controller
	option *ControllerOptions
}

// NewRoomController creates a room controller.
func NewRoomController(service *goa.Service, option *ControllerOptions) *RoomController {
	return &RoomController{
		Controller: service.NewController("RoomController"),
		option:     option,
	}
}

// List runs the list action.
func (c *RoomController) List(ctx *app.ListRoomContext) error {

	db := c.option.db
	// connections := c.option.connections

	res := app.RoomCollection{}
	rooms, err := models.AllRooms(db, 100)
	if err != nil {
		return err
	}
	for _, room := range rooms {
		res = append(res, ToRoomMedia(room))
	}
	return ctx.OK(res)
}

// Post runs the post action.
func (c *RoomController) Post(ctx *app.PostRoomContext) error {
	db := c.option.db
	// connections := c.option.connections

	token := jwt.ContextJWT(ctx)
	if token == nil {
		return fmt.Errorf("JWT token is missing from context") // internal error
	}
	claims := token.Claims.(jwtgo.MapClaims)

	loginfo(ctx, "func (c *RoomController) Post(ctx *app.PostRoomContext) ",
		"claims", claims)

	// Use the claims to authorize
	// subject := claims["sub"]
	// if subject != "subject" {
	// 	// A real app would probably use an "Unauthorized" response here
	// 	// res := &app.Success{OK: false}
	// 	// return ctx.OK(res)
	// }

	room := models.Room{
		Name:        ctx.Payload.Name,
		Description: ctx.Payload.Description,
		Created:     time.Now(),
	}
	err := room.Insert(db)
	if err != nil {
		return err
	}
	str := app.RoomHref(room.ID)
	loginfo(ctx, "str,", str)
	ctx.ResponseData.Header().Set("Location", app.RoomHref(room.ID))
	return ctx.Created()
}

// Show runs the show action.
func (c *RoomController) Show(ctx *app.ShowRoomContext) error {
	db := c.option.db
	// connections := c.option.connections
	// if room, ok := c.db.GetRoom(ctx.RoomID); ok {
	room, err := models.RoomByID(db, ctx.RoomID)
	if err != nil {
		return err
	}
	if room == nil {
		return ctx.NotFound()
	}
	res := ToRoomMedia(room)
	return ctx.OK(res)
}

// Watch watches the message with the given id.
func (c *RoomController) Watch(ctx *app.WatchRoomContext) error {
	Watcher(ctx.RoomID, c, ctx).ServeHTTP(ctx.ResponseWriter, ctx.Request)
	return nil
}

// Roomの変更通知の送信
func Watcher(roomID int, c *RoomController, ctx *app.WatchRoomContext) websocket.Handler {
	// db := c.option.db
	connections := c.option.connections
	return func(ws *websocket.Conn) {
		ch := make(chan struct{})
		connections.apendConn(roomID, ch)
		for {
			<-ch
			_, err := ws.Write([]byte(fmt.Sprintf("Room: %d", roomID)))
			if err != nil {
				break
			}
		}
		connections.removeConn(roomID, ch)
	}
}
