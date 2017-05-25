package controllers

import (
	"github.com/goadesign/goa"
	"github.com/m0a-mystudy/goa-chat/app"
)

// AccountController implements the account resource.
type AccountController struct {
	*goa.Controller
	option *ControllerOptions
}

// NewAccountController creates a account controller.
func NewAccountController(service *goa.Service, option *ControllerOptions) *AccountController {
	return &AccountController{
		Controller: service.NewController("AccountController"),
		option:     option,
	}
}

// List runs the list action.
func (c *AccountController) List(ctx *app.ListAccountContext) error {
	// AccountController_List: start_implement

	// Put your logic here

	// AccountController_List: end_implement
	res := app.AccountCollection{}
	return ctx.OK(res)
}

// Post runs the post action.
func (c *AccountController) Post(ctx *app.PostAccountContext) error {
	// AccountController_Post: start_implement

	// Put your logic here

	// AccountController_Post: end_implement
	return nil
}

// Show runs the show action.
func (c *AccountController) Show(ctx *app.ShowAccountContext) error {
	// AccountController_Show: start_implement

	// Put your logic here

	// AccountController_Show: end_implement
	res := &app.Account{}
	return ctx.OK(res)
}
