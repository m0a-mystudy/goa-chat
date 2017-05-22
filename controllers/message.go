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
	db          *sql.DB
	connections *WsConnections
}

// NewMessageController creates a message controller.
func NewMessageController(service *goa.Service, db *sql.DB, wsc *WsConnections) *MessageController {
	return &MessageController{
		Controller:  service.NewController("MessageController"),
		db:          db,
		connections: wsc,
	}
}

// List runs the list action.
func (c *MessageController) List(ctx *app.ListMessageContext) error {
	res := app.MessageCollection{}

	option := models.MessageParamOption{
		RoomID:          ctx.RoomID,
		Limit:           100,
		Offset:          0,
		OrderByPostDate: true,
	}

	if ctx.Limit != nil {
		option.Limit = *ctx.Limit
	}
	if ctx.Offset != nil {
		option.Offset = *ctx.Offset
	}

	messages, err := models.MessagesByOption(c.db, option)
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
	c.connections.updateRoom(ctx.RoomID)
	ctx.ResponseData.Header().Set("Location", app.MessageHref(ctx.RoomID, m.ID))
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
