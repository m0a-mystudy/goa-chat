goaとxoでchat apiを作ってみた。




こんな感じのurlでwebAPIを作りたいかなと漠然と考えました

| methods | エンドポイント | 目的|
|--------|--------|-----|
| get,post |http://chat/api/rooms |チャットルームの作成と一覧の取得 |
| get,post |http://chat/api/rooms/3/messages |room 3へのチャットメッセージの作成と一覧の取得 |


これを念頭にDSLを書きます

あとデータベースはmysqlを使います。xoでdbからコード生成を行う方針で行きます。


# 初期設定

goaのインストール

```bash
go get -u -v  github.com/goadesign/goa/..
```

``GOPATH``上にプロジェクトを作ります

```bash:自分の場合
mkdir -p $GOPATH/src/github.com/m0a-mystudy/goa-chat
cd $GOPATH/src/github.com/m0a-mystudy/goa-chat
```

[goa-celler](https://github.com/goadesign/goa-cellar)を参考にMakefileを作ります。

```Makefile

#! /usr/bin/make
#
# Makefile for goa chat
#
# Targets:
# - clean     delete all generated files
# - generate  (re)generate all goagen-generated files.
# - build     compile executable
#
# Meta targets:
# - all is the default target, it runs all the targets in the order above.
#

all: depend clean generate build

depend:
	@glide install

clean:
	@rm -rf app
	@rm -rf client
	@rm -rf tool
	@rm -rf public/swagger
	@rm -rf public/schema
	@rm -rf public/js
	@rm -f todo

bootstrap:
	@goagen main    -d github.com/m0a-mystudy/goa-chat/design -o controllers

generate:
	@goagen app     -d github.com/m0a-mystudy/goa-chat/design
	@goagen swagger -d github.com/m0a-mystudy/goa-chat/design -o public
	@goagen schema  -d github.com/m0a-mystudy/goa-chat/design -o public
	@goagen client  -d github.com/m0a-mystudy/goa-chat/design
	@goagen js      -d github.com/m0a-mystudy/goa-chat/design -o public

build:
	@go build -o chat


```




# DSLの記述



```go:design/design.go
package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("Chat API", func() {
	Title("goa study chat") // Documentation title
	Description("goa study chat api")
	Host("localhost:8080")
	Scheme("http")
	BasePath("/api")
	
	// クライアントを別に作ろうとしていたのでとりあえずクロスサイト可能にしておく
	Origin("http://localhost:3000", func() {
		Methods("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS")
		Headers("Origin", "X-Requested-With", "Content-Type", "Accept")
	})
})

var _ = Resource("room", func() {
	DefaultMedia(Room)
	BasePath("/rooms")
	Action("list", func() {
		Routing(GET(""))
		Description("Retrieve all rooms.")
		Response(OK, CollectionOf(Room))
		Response(NotFound)
	})

	Action("show", func() {
		Routing(
			GET("/:roomID"),
		)
		Description("Retrieve room with given id")
		Params(func() {
			Param("roomID", Integer)
		})
		Response(OK)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})

	Action("post", func() {
		Routing(POST(""))
		Description("Create new Room")
		Payload(RoomPayload)
		Response(Created, Room)
		Response(BadRequest)
	})

})

var _ = Resource("message", func() {
	DefaultMedia(Message)
	BasePath("messages")
	Parent("room")
	Action("list", func() {
		Routing(GET(""))
		Description("Retrieve all messages.")
		Response(OK, CollectionOf(Message))
		Response(NotFound)
	})
	Action("post", func() {
		Routing(POST(""))
		Description("Create new message")
		Payload(MessagePayload)
		Response(Created, Message)
		Response(BadRequest)
	})

})

var Message = MediaType("application/vnd.message+json", func() {
	Description("A Message")
	Reference(MessagePayload)
	Attributes(func() {
		Attribute("accountID")
		Attribute("body")
		Attribute("postDate")
		Required("accountID", "body", "postDate")
	})
	View("default", func() {
		Attribute("accountID")
		Attribute("body")
		Attribute("postDate")
	})
})
var MessagePayload = Type("MessagePayload", func() {

	Attribute("accountID", Integer, func() {
		Example(1)
	})
	Attribute("body", func() {
		MinLength(1)
		MaxLength(400)
		Example("this is chat message")
	})
	Attribute("postDate", DateTime, func() {
		Default("1978-06-30T10:00:00+09:00")
	})

	Required("accountID", "body", "postDate")
})

var Room = MediaType("application/vnd.room+json", func() {
	Description("A room")
	Reference(RoomPayload)
	Attributes(func() {
		Attribute("id")
		Attribute("name")
		Attribute("description")
		Attribute("created")
		Required("name", "description")
	})

	View("default", func() {
		Attribute("id")
		Attribute("name")
		Attribute("description")
		Attribute("created")

	})
})

var RoomPayload = Type("RoomPayload", func() {
	Attribute("id", Integer, "ID of room")
	Attribute("name", String, "Name of room", func() {
		Example("room001")
	})
	Attribute("description", String, "description of room", func() {
		Example("room description")
		MaxLength(400)
	})
	Attribute("created", DateTime, "Date of creation")
	Required("name", "description")
})

```


上記はまだpayloadの作りとかよくわかってないのでちょっと仮な感じです。
ちょっとずつこちらを修正していきます。

## 詰まったところ


### CORS設定
まずクライアントは``create-react-app``を使って作ろうとしてまして、何も考えないと別ドメイン扱いになります。
それでもアクセス可能にするためにCORSの設定が必要でしたがどこにも情報がなくて苦労しました。
(``Content-Type is not allowed by Access-Control-Allow-Headers``とか出る)
以下のように設定したら上手く動きました。

```go
    Origin("http://localhost:3000", func() {
        Methods("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS")
        Headers("Origin", "X-Requested-With", "Content-Type", "Accept")
    })
```

参考：http://stackoverflow.com/questions/5027705/error-in-chrome-content-type-is-not-allowed-by-access-control-allow-headers


### 親子関係の設定

| methods | エンドポイント | 目的|
|--------|--------|-----|
| get,post |http://chat/api/rooms/3/messages |room 3へのチャットメッセージの作成と一覧の取得 |

上記のように２つのリソース roomとmessageには親子関係になっているので
それを素直に設定するとこんな感じになります


```go
var _ = Resource("room", func() {
	DefaultMedia(Room)
	BasePath("/rooms")

	Action("list", func() {
		Routing(GET(""))
		//省略
	})

	Action("post", func() {
		Routing(POST(""))
		//省略
	})

})

var _ = Resource("message", func() {
	DefaultMedia(Message)
	BasePath("messages")
	Parent("room")
	
	Action("list", func() {
		Routing(GET(""))
		//省略
	})
	Action("post", func() {
		Routing(POST(""))
		//省略
	})

})
```

ところがこのままdslからコード生成を行うと以下のようなエラーが出ます

```
exit status 1
resource "message": Parent resource "room" has no canonical action
link "room" of type "Message": Link name must match one of the parent media type attribute names
make: *** [bootstrap] Error 1
```

このエラーメッセージに悩みました。
対応方法は``show``actionを親側に定義します。(canonical action)


```go:roomリソースに以下のactionを追加
	Action("show", func() {
		Routing(
			GET("/:roomID"),
		)
		Description("Retrieve room with given id")
		Params(func() {
			Param("roomID", Integer)
		})
		Response(OK)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
	})
```


参考：https://goa.design/reference/goa/design/apidsl/#func-resource-a-name-apidsl-resource-a

(ちなみに``show``という名前が嫌ならCanonicalActionNameを使って変更できる)


# 実装

``make bootstrap``を実行してコード生成を行います。
修正が必要なファイルは``controllers``に作られますが、main.goだけは手動で直下に移動しておきます。

## database側を作る

goaをつかってdbのスキーマを作るのもいいのですが今回は
mysql側でsqlを書いてxoでgoの構造体を作るとう言う方針で行きます。

xoのインストール

```bash
$ go get -u -v github.com/knq/xo...
```

mysql側でスキーマを作りました。ダンプは以下となります。

```
--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `messages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `room_id` int(11) NOT NULL,
  `account_id` int(11) NOT NULL,
  `body` varchar(400) NOT NULL,
  `postDate` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `room_id_idx` (`room_id`),
  CONSTRAINT `room_id` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `rooms`
--

DROP TABLE IF EXISTS `rooms`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rooms` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(256) NOT NULL,
  `description` varchar(400) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name_UNIQUE` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;


```

mysqlへ流し込んでおきます。(方法は省略)
db_nameは``goa_chat``とします。

以下のコマンドでdbスキーマからgo codeを生成します。

```bash
$ mkdir -p models
$ xo  mysql://<usrname>:<pass>@localhost/goa_chat -o models
$ ls -l models       
total 64
-rw-r--r--  1 m0a  staff  4199  5 17 13:25 message.xo.go
-rw-r--r--  1 m0a  staff  3581  5 17 13:25 room.xo.go
-rw-r--r--  1 m0a  staff  2128  5 17 10:56 xo_db.xo.go
```

生成されたコードはこんな感じです

```go:models/room.xo.go
// Package models contains the types for schema 'goa_chat'.
package models

// GENERATED BY XO. DO NOT EDIT.

import (
	"errors"
	"time"
)

// Room represents a row from 'goa_chat.rooms'.
type Room struct {
	ID          int       `json:"id"`          // id
	Name        string    `json:"name"`        // name
	Description string    `json:"description"` // description
	Created     time.Time `json:"created"`     // created

	// xo fields
	_exists, _deleted bool
}

// Exists determines if the Room exists in the database.
func (r *Room) Exists() bool {
	return r._exists
}

// Deleted provides information if the Room has been deleted from the database.
func (r *Room) Deleted() bool {
	return r._deleted
}

// Insert inserts the Room to the database.
func (r *Room) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if r._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO goa_chat.rooms (` +
		`name, description, created` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, r.Name, r.Description, r.Created)
	res, err := db.Exec(sqlstr, r.Name, r.Description, r.Created)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	r.ID = int(id)
	r._exists = true

	return nil
}

//
//     コード省略
//

// RoomByID retrieves a row from 'goa_chat.rooms' as a Room.
//
// Generated from index 'rooms_id_pkey'.
func RoomByID(db XODB, id int) (*Room, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`id, name, description, created ` +
		`FROM goa_chat.rooms ` +
		`WHERE id = ?`

	// run query
	XOLog(sqlstr, id)
	r := Room{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, id).Scan(&r.ID, &r.Name, &r.Description, &r.Created)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

```

中身を見ると全体を取得する関数が定義されてないのでroomの分だけ作ります

```go:models/room.go
package models

func AllRooms(db XODB, limit int) ([]*Room, error) {

	// sql query
	const sqlstr = `SELECT ` +
		`id, name, description, created ` +
		`FROM goa_chat.rooms LIMIT ?`
	// run query
	XOLog(sqlstr, limit)
	q, err := db.Query(sqlstr, limit)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	var res []*Room
	for q.Next() {
		r := Room{}
		err = q.Scan(&r.ID, &r.Name, &r.Description, &r.Created)
		if err != nil {
			return nil, err
		}
		res = append(res, &r)
	}
	return res, nil
}

```

## goaのコードとdbの紐付け

goa-cellerを真似てcontrollerにdbをつけておきます。

``controllers``の``room.go``,``message.go``を修正します


```diff:controllers/room.go
// RoomController implements the room resource.
type RoomController struct {
	*goa.Controller
+	db *sql.DB
}

// NewRoomController creates a room controller.
- func NewRoomController(service *goa.Service) *RoomController {
+ func NewRoomController(service *goa.Service, db *sql.DB) *RoomController {
	return &RoomController{
		Controller: service.NewController("RoomController"),
+		db:         db,
	}
}
```

``message.go``は省略します

``main.go``にて実際にdbを接続するコードを追加します
 
```diff:main.go
package main

import (
+	"database/sql"

	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/controllers"

+	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Create service
	service := goa.New("Chat API")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

+	db, err := sql.Open("mysql", "user:password@/goa_chat?parseTime=true")
+	if err != nil {
+		service.LogError("startup", "err", err)
+	}
	// Mount "message" controller
-	c := controllers.NewMessageController(service)
+	c := controllers.NewMessageController(service, db)
	app.MountMessageController(service, c)
	// Mount "room" controller
-	c2 := controllers.NewRoomController(service)
+	c2 := controllers.NewRoomController(service, db)
	app.MountRoomController(service, c2)

	// Start service
	if err := service.ListenAndServe(":8080"); err != nil {
		service.LogError("startup", "err", err)
	}
}

```

## controllers/room.go実装

先ずはmodel側の構造体とgoa側の構造体の変換用の処理関数を作っておきます

```go
func ToRoomMedia(room *models.Room) *app.Room {
	ret := app.Room{
		ID:          &room.ID,
		Description: room.Description,
		Name:        room.Name,
		Created:     &room.Created,
	}
	return &ret
}

```

DSLでActionを定義した分だけメソッドができているので中味を実装していきます。

```go
package controllers

import (
	"database/sql"
	"time"

	"github.com/goadesign/goa"
	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/models"
)


//
//     コード省略
//

// List runs the list action.
func (c *RoomController) List(ctx *app.ListRoomContext) error {
	res := app.RoomCollection{}
	rooms, err := models.AllRooms(c.db, 100) //とりあえず100件固定で
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

```

ctxからroomIDが取得できたり、Actionのルーティング設定から必要なパラメータをctxから取得できます。
便利。フレームワークによってはinterface{}型だったりするんですがコード生成なのでちゃんとroomIDはint型になっているのが素敵です。

## controllers/message.go実装

こちらも基本的に同じです。

```go
package controllers

import (
	"database/sql"
	"time"

	"github.com/goadesign/goa"
	"github.com/m0a-mystudy/goa-chat/app"
	"github.com/m0a-mystudy/goa-chat/models"
)

//
//     コード省略
//


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
		return ctx.BadRequest()
	}

	return ctx.Created(ToMessageMedia(&m))
}

```

xoもスキーマから読み取って``MessagesByRoomID``を作ってくれるのが素敵です。


# クライアント側の実装

自分はtypeScriptが好きなのでreact+typescriptな環境を作ります。

```bash
$ create-react-app --scripts-version=react-scripts-ts  goa-chat-client
$ cd goa-chat-client
$ yarn 
$ yarn start
```

``localhost:3000``でクライアントが立ち上がります。

swaggerからclient apiを作るための環境をインストールします

```bash
$ brew install swagger-codegen
```

``goa-chat-client/src/comm``にclient apiのコードを出力します

```bash
$ cd goa-chat-client/src
$ mkdir comm
$ cd comm
$ swagger-codegen generate -l typescript-fetch -i ../../../public/swagger/swagger.json
```

client apiをみながら必要なパッケージを追加しておきます

```
$ cd goa-chat-client
$ yarn add isomorphic-fetch core-js --save
$ yarn add @types/isomorphic-fetch @types/core-js --save-dev
```

それだけでは動かずエラー内容を見ながら
`` import * as querystring from "querystring";``の行も削除しました。


実際に使う場合は非常に簡単でした

```typescript
import * as comm from './comm/api';

let messageAPI = new comm.MessageApi();
const roomID = 10;
let messages = await messageAPI.messageList({ roomID });
```
と言った具合に取得できます。(擬似コードです。実際には動かないです)

あとはreactで画面を作っていけばいいんですが
ここで力尽きました。


# 最後に

swagger.yamlが作られているので
https://editor.swagger.io/ から対応するcurlコマンドを作ってくれたり
typescriptのクライアントapiを自動生成してくれるのが素敵です。


とりあえずここまでのコードは以下においておきます。
https://github.com/m0a-mystudy/goa-chat

goaはwebsocketもサポートしてるのでそれにも対応させてみたいです。








