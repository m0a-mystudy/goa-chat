package controllers

import "context"

type Comm chan struct{}

// WsConnections is connection pool websocket conn
type WsConnections struct {
	connections map[int][]Comm
	ctx         context.Context
}

// NewConnections create WsConnections
func NewConnections(ctx context.Context) *WsConnections {
	return &WsConnections{
		connections: map[int][]Comm{},
		ctx:         ctx,
	}
}

func (l *WsConnections) apendConn(roomID int, comm Comm) {
	list := l.connections[roomID]
	if list == nil {
		list = []Comm{}
	}
	list = append(list, comm)
	l.connections[roomID] = list
	loginfo(l.ctx, "apendConn", "list", list)
}

func (l *WsConnections) removeConn(roomID int, comm Comm) {
	list := l.connections[roomID]
	if list == nil {
		list = []Comm{}
	}
	newList := []Comm{}
	for _, c := range list {
		if c == comm {
			continue
		}
		newList = append(newList, c)
	}

	l.connections[roomID] = newList
	loginfo(l.ctx, "removeConn", "list", newList)
}

func (l *WsConnections) updateRoom(roomID int) {
	loginfo(l.ctx, "updateRoom", "roomID", roomID)
	comms, ok := l.connections[roomID]
	if !ok {
		return
	}
	loginfo(l.ctx, "updateRoom", "comms", comms)
	for _, comm := range comms {
		comm <- struct{}{}
	}
}
