package service

import (
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketService struct {
	UserRepo IUserRepo
}

type Chat struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Msg       string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Client struct {
	Conn     *websocket.Conn
	Username string
}

type Message struct {
	Type string `json:"type"`
	User string `json:"user,omitempty"`
	Chat Chat   `json:"chat,omitempty"`
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *Chat)

// We'll need to define an Upgrader
// this will require a Read and Write buffer size
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// We'll need to check the origin of our connection
	// this will allow us to make requests from our React
	// development server to here.
	// For now, we'll do no checking and just allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWebSocketService(userRepo IUserRepo) *WebSocketService {
	return &WebSocketService{UserRepo: userRepo}
}

// define our WebSocket endpoint
func (s *WebSocketService) ServeWs(ctx context.Context, username string, w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	go Broadcaster(ctx)

	fmt.Println(r.Host, r.URL.Query())

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{Conn: ws, Username: username}
	// register client
	clients[client] = true
	fmt.Println("clients", len(clients), clients, ws.RemoteAddr())

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	Receiver(ctx, client)

	fmt.Println("exiting", ws.RemoteAddr().String())
	delete(clients, client)
}

// define a receiver which will listen for
// new messages being sent to our WebSocket
// endpoint
func Receiver(ctx context.Context, client *Client) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	for {
		// read in a message
		// readMessage returns messageType, message, err
		// messageType: 1-> Text Message, 2 -> Binary Message
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		m := &Message{}

		err = json.Unmarshal(p, m)
		if err != nil {
			log.Println("error while unmarshaling chat", err)
			continue
		}

		fmt.Println("host", client.Conn.RemoteAddr())
		if m.Type == "bootup" {
			// do mapping on bootup
			client.Username = m.User
			fmt.Println("client successfully mapped", &client, client, client.Username)
		} else {
			fmt.Println("received message", m.Type, m.Chat)
			c := m.Chat
			c.Timestamp = time.Now().Unix()

			// save in redis
			// id, err := redisrepo.CreateChat(&c)
			id := "this is chat id"
			// if err != nil {
			// 	log.Println("error while saving chat in redis", err)
			// 	return
			// }

			c.ID = id
			broadcast <- &c
		}
	}
}

func Broadcaster(ctx context.Context) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	for {
		fmt.Println("Printed every second")
		time.Sleep(time.Second)
		message := <-broadcast
		// send to every client that is currently connected
		fmt.Println("new message", message)

		for client := range clients {
			// send message only to involved users
			fmt.Println("username:", client.Username,
				"from:", message.From,
				"to:", message.To)

			if client.Username == message.From || client.Username == message.To {
				err := client.Conn.WriteJSON(message)
				if err != nil {
					log.Printf("Websocket error: %s", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}

// func setupRoutes() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintf(w, "Simple Server1")
// 	})
// 	// map our `/ws` endpoint to the `serveWs` function
// 	http.HandleFunc("/ws", ServeWs)
// }

// func StartWebsocketServer() {
// 	// redisClient := redisrepo.InitialiseRedis()
// 	// defer redisClient.Close()

// 	go broadcaster()
// 	setupRoutes()
// 	http.ListenAndServe(":8080", nil)
// }
