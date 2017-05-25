package controllers

import (
	"crypto/rsa"
	"database/sql"
)

// ControllerOptions is Controller option
type ControllerOptions struct {
	db          *sql.DB
	connections *WsConnections
	privateKey  *rsa.PrivateKey
}

// NewOption is create ControllerOptions
func NewOption(
	db *sql.DB,
	wsc *WsConnections,
	privateKey *rsa.PrivateKey,
) *ControllerOptions {

	return &ControllerOptions{
		db:          db,
		connections: wsc,
		privateKey:  privateKey,
	}
}


