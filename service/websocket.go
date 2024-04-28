package service

import (
	"app/entity"
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
	UserRepo   IUserRepo
	SocketRepo ISocketRepo
}

type Client struct {
	Conn   *websocket.Conn
	UserID string
}

var clients = make(map[*Client]bool)
var broadcast = make(chan *entity.Chat)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewWebSocketService(userRepo IUserRepo, socketRepo ISocketRepo) *WebSocketService {
	return &WebSocketService{UserRepo: userRepo, SocketRepo: socketRepo}
}

func (s *WebSocketService) ServeWs(ctx context.Context, userID string, w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	go s.Broadcaster(ctx)

	fmt.Println(r.Host, r.URL.Query())

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	client := &Client{Conn: ws, UserID: userID}
	clients[client] = true
	fmt.Println("clients", len(clients), clients, ws.RemoteAddr())
	s.Receiver(ctx, client)

	fmt.Println("exiting", ws.RemoteAddr().String())
	delete(clients, client)
}

func (s *WebSocketService) Receiver(ctx context.Context, client *Client) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	for {
		_, p, err := client.Conn.ReadMessage()
		if err != nil {
			panic(err)
		}

		m := &entity.Message{}

		err = json.Unmarshal(p, m)
		if err != nil {
			panic(err)
		}

		fmt.Println("host", client.Conn.RemoteAddr())

		if m.Type == "newChat" {
			newRepo := &entity.Conversation{
				ID:       utils.NewID(),
				Name:     "new_chat",
				ListUser: m.ListUserInNewChat,
				Chat:     []entity.Chat{m.Chat},
			}
			err = s.SocketRepo.NewConversation(ctx, newRepo)
			if err != nil {
				panic(err)
			}
			client.UserID = m.Chat.FromUserId
		} else {
			fmt.Println("received message", m.Type, m.Chat)
			c := m.Chat
			c.Timestamp = time.Now()

			err = s.SocketRepo.AddNewChatToConversation(ctx, &m.Chat)
			if err != nil {
				panic(err)
			}

			c.ID = utils.NewID()
			broadcast <- &c
		}
	}
}

func (s *WebSocketService) Broadcaster(ctx context.Context) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()
	for {
		message := <-broadcast
		fmt.Println("new message", message)

		for client := range clients {
			fmt.Println("userID:", client.UserID,
				"from:", message.FromUserId,
				"to:", message.ToConversationId)

			listUser, err := s.SocketRepo.GetListUserInConversation(ctx, message.ToConversationId)
			if err != nil {
				panic(err)
			}
			if client.UserID == message.FromUserId || utils.ContainsString(listUser, client.UserID) {
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