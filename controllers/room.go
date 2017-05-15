package controllers

import (
	"github.com/goadesign/goa"
	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/store"
)

// ToRoomMedia convert tool
func ToRoomMedia(room store.RoomModel) *app.Room {
	ret := app.Room(room)
	return &ret
}

// RoomController implements the room resource.
type RoomController struct {
	*goa.Controller
	db *store.DB
}

// NewRoomController creates a room controller.
func NewRoomController(service *goa.Service, db *store.DB) *RoomController {
	return &RoomController{
		Controller: service.NewController("RoomController"),
		db:         db,
	}
}

// List runs the list action.
func (c *RoomController) List(ctx *app.ListRoomContext) error {
	res := app.RoomCollection{}
	list := c.db.GetRooms()
	for _, room := range list {
		res = append(res, ToRoomMedia(room))
	}
	return ctx.OK(res)
}

// Post runs the post action.
func (c *RoomController) Post(ctx *app.PostRoomContext) error {
	model := c.db.NewRoom()
	saveModel := store.RoomModel(*ctx.Payload)
	saveModel.ID = model.ID
	c.db.SaveRoom(saveModel)
	return nil
}

// Show runs the show action.
func (c *RoomController) Show(ctx *app.ShowRoomContext) error {
	if room, ok := c.db.GetRoom(ctx.RoomID); ok {
		res := ToRoomMedia(room)
		return ctx.OK(res)
	}
	return ctx.NotFound()
}
