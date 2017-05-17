package models

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDatabase(t *testing.T) {
	db, err := sql.Open("mysql", "root:sanset6@/goa_chat?parseTime=true")
	if err != nil {
		t.Error(err)
	}

	id, err := NewRoom(db, "myfirstroom02", "my first room")
	if err != nil {
		t.Error(err)
	}

	room, err := RoomByID(db, id)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(room)
}
