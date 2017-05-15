package store

import (
	"fmt"
	"sync"

	"github.com/m0a-mystudy/goa-chat/app"
)

// MessageModel Payloadから作成
type MessageModel app.MessagePayload
type RoomModel app.RoomPayload

// DB 本体
type DB struct {
	sync.Mutex
	maxRoomID int
	rooms     map[int]*RoomModel
	messages  map[int][]*MessageModel
}

// CreateTodoPayload
// // TodoModel Todo のモデル
// type TodoModel struct {
// 	ID       int
// 	Title    string
// 	Created  time.Time
// 	Modified time.Time
// }

// NewDB は新規DB作成
func NewDB() *DB {
	return &DB{
		rooms:    map[int]*RoomModel{},
		messages: map[int][]*MessageModel{},
	}
}

func (db *DB) GetRooms() []RoomModel {
	db.Lock()
	defer db.Unlock()
	fmt.Printf("db.rooms =%#v", db.rooms)
	list := make([]RoomModel, len(db.rooms))
	fmt.Printf("list =%#v", list)
	for i, room := range db.rooms {
		list[i] = *room
	}
	return list
}

func (db *DB) GetRoom(room int) (model RoomModel, ok bool) {
	db.Lock()
	defer db.Unlock()
	if model, ok := db.rooms[room]; ok {
		return *model, true
	}
	return model, false
}

// func (db *DB) GetMessage(room, id int) (model MessageModel, ok bool) {
// 	db.Lock()
// 	defer db.Unlock()
// 	messages, found := db.messages[room]
// 	if !found {
// 		return
// 	}

// 	for _, m := range messages {
// 		if *m.ID == id {
// 			model = *m
// 			ok = true
// 			return
// 		}
// 	}
// 	return
// }

func (db *DB) GetMessages(room int) ([]MessageModel, error) {
	db.Lock()
	defer db.Unlock()
	_, ok := db.rooms[room]
	if !ok {
		return nil, fmt.Errorf("unkown room %d", room)
	}
	messages, ok := db.messages[room]
	if !ok {
		db.messages[room] = []*MessageModel{}
		messages = []*MessageModel{}
	}
	list := make([]MessageModel, len(messages))
	for i, m := range messages {
		list[i] = *m
	}
	return list, nil
}

// NewMessage は新規作成
func (db *DB) NewMessage() (model MessageModel) {
	model = MessageModel{}
	return
}

// SaveMessage によるセーブ
func (db *DB) SaveMessage(room int, model MessageModel) error {
	db.Lock()
	defer db.Unlock()
	messages, ok := db.messages[room]
	if !ok {
		return fmt.Errorf("unkown room %d", room)
	}
	messages = append(messages, &model)
	db.messages[room] = messages
	return nil
}

// NewRoom は新規作成
func (db *DB) NewRoom() (model RoomModel) {
	id := db.maxRoomID
	db.maxRoomID++

	model = RoomModel{
		ID: &id,
	}
	return
}

func (db *DB) SaveRoom(model RoomModel) {
	db.rooms[*model.ID] = &model
}
