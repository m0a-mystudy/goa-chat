package controllers

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/net/websocket"

	"github.com/goadesign/goa"
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
	db          *sql.DB
	connections *WsConnections
}

// NewRoomController creates a room controller.
func NewRoomController(service *goa.Service, db *sql.DB, wsc *WsConnections) *RoomController {
	return &RoomController{
		Controller:  service.NewController("RoomController"),
		db:          db,
		connections: wsc,
	}
}

// List runs the list action.
func (c *RoomController) List(ctx *app.ListRoomContext) error {
	res := app.RoomCollection{}
	rooms, err := models.AllRooms(c.db, 100)
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
	room := models.Room{
		Name:        ctx.Payload.Name,
		Description: ctx.Payload.Description,
		Created:     time.Now(),
	}
	err := room.Insert(c.db)
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
	// if room, ok := c.db.GetRoom(ctx.RoomID); ok {
	room, err := models.RoomByID(c.db, ctx.RoomID)
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
	return func(ws *websocket.Conn) {
		ch := make(chan struct{})
		c.connections.apendConn(roomID, ch)
		for {
			<-ch
			_, err := ws.Write([]byte(fmt.Sprintf("Room: %d", roomID)))
			if err != nil {
				break
			}
		}
		c.connections.removeConn(roomID, ch)
	}
}
