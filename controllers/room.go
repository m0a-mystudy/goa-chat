package controllers

import (
	"database/sql"
	"time"

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
	db *sql.DB
}

// NewRoomController creates a room controller.
func NewRoomController(service *goa.Service, db *sql.DB) *RoomController {
	return &RoomController{
		Controller: service.NewController("RoomController"),
		db:         db,
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
	return ctx.Created(ToRoomMedia(&room))
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
