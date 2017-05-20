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
	Origin("http://localhost:3000", func() {
		Methods("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS")
		Headers("Origin", "X-Requested-With", "Content-Type", "Accept")
	})
	ResponseTemplate(Created, func(pattern string) {
		Description("Resource created")
		Status(201)
		Headers(func() {
			Header("Location", String, "href to created resource", func() {
				Pattern(pattern)
			})
		})
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
		Response(Created, "/rooms/[0-9]+")
		Response(BadRequest)
	})

	Action("watch", func() {
		Routing(
			GET("/:roomID/watch"),
		)
		Scheme("ws")
		Description("Retrieve room with given id")
		Params(func() {
			Param("roomID", Integer)
		})
		Response(SwitchingProtocols)
		Response(BadRequest, ErrorMedia)
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
		Response(Created, "^/rooms/[0-9]+/messages/[0-9]+$")
		Response(BadRequest)
	})

	Action("show", func() {
		Routing(
			GET("/:messageID"),
		)
		Description("Retrieve message with given id")
		Params(func() {
			Param("messageID", Integer)
		})
		Response(OK)
		Response(NotFound)
		Response(BadRequest, ErrorMedia)
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
