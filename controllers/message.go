package controllers

import (
	"github.com/goadesign/goa"
	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/store"
)

// ToMessageMedia convert tool
func ToMessageMedia(model store.MessageModel) *app.Message {
	ret := app.Message(model)
	return &ret
}

// MessageController implements the message resource.
type MessageController struct {
	*goa.Controller
	db *store.DB
}

// NewMessageController creates a message controller.
func NewMessageController(service *goa.Service, db *store.DB) *MessageController {
	return &MessageController{
		Controller: service.NewController("MessageController"),
		db:         db,
	}
}

// List runs the list action.
func (c *MessageController) List(ctx *app.ListMessageContext) error {
	res := app.MessageCollection{}
	messages, err := c.db.GetMessages(ctx.RoomID)
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
	roomID := ctx.RoomID
	model := store.MessageModel(*ctx.Payload)
	if c.db.SaveMessage(roomID, model) != nil {
		ctx.Created(ToMessageMedia(model))
	}

	return ctx.BadRequest()
}
