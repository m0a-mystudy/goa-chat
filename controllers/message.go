package controllers

import (
	"database/sql"
	"time"

	"github.com/goadesign/goa"
	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/models"
)

// ToMessageMedia convert tool
func ToMessageMedia(model *models.Message) *app.Message {
	ret := app.Message{
		Body:      model.Body,
		AccountID: model.AccountID,
		PostDate:  model.Postdate,
	}
	return &ret
}

// MessageController implements the message resource.
type MessageController struct {
	*goa.Controller
	db *sql.DB
}

// NewMessageController creates a message controller.
func NewMessageController(service *goa.Service, db *sql.DB) *MessageController {
	return &MessageController{
		Controller: service.NewController("MessageController"),
		db:         db,
	}
}

// List runs the list action.
func (c *MessageController) List(ctx *app.ListMessageContext) error {
	res := app.MessageCollection{}

	messages, err := models.MessagesByRoomID(c.db, ctx.RoomID)
	if err != nil {
		return err
	}
	for _, m := range messages {
		res = append(res, ToMessageMedia(m))
	}
	return ctx.OK(res)
}

// Post runs the post action.
func (c *MessageController) Post(ctx *app.PostMessageContext) error {
	m := models.Message{
		RoomID:    ctx.RoomID,
		AccountID: ctx.Payload.AccountID,
		Body:      ctx.Payload.Body,
		Postdate:  time.Now(),
	}

	err := m.Insert(c.db)
	if err != nil {
		//return err
		return ctx.BadRequest()
	}

	ctx.Response.Header.Set("Location", app.MessageHref(ctx.RoomID, m.ID))
	return ctx.Created()
}

// Show runs the show action.
func (c *MessageController) Show(ctx *app.ShowMessageContext) error {
	// if room, ok := c.db.GetRoom(ctx.RoomID); ok {
	message, err := models.MessageByID(c.db, ctx.MessageID)
	if err != nil {
		return err
	}
	if message == nil {
		return ctx.NotFound()
	}
	res := ToMessageMedia(message)
	return ctx.OK(res)
}
