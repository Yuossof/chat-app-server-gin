package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Yuossof/messaging-app-server/database"
	"github.com/Yuossof/messaging-app-server/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocket Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Store active WebSocket connections for each user
var clients = make(map[uuid.UUID]*websocket.Conn) // userID -> WebSocket connection
var mu sync.Mutex

// Channel for message processing
var messageChannel = make(chan models.Message, 1000)
var wg sync.WaitGroup

func init() {
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		go processMessages(i)
	}
}

func processMessages(workerID int) {
	defer wg.Done()

	for msg := range messageChannel {
		result := database.DB.Create(&msg)
		if result.Error != nil {
			log.Println("❌ Error saving message:", result.Error)
		} else {
			log.Printf("✅ Worker %d saved message: %s", workerID, msg.Content)
		}
	}
}

// WebSocket Handler
func WebSocketHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("❌ Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Get user_id from query parameters
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		log.Println("❌ Invalid user ID:", err)
		return
	}

	// Store user connection
	mu.Lock()
	clients[userID] = conn
	mu.Unlock()

	log.Printf("✅ User %s connected to WebSocket\n", userID)

	// Listen for messages from the client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("❌ Connection lost with user:", userID)
			mu.Lock()
			delete(clients, userID)
			mu.Unlock()
			break
		}

		// Decode the received message
		var payload models.Message
		err = json.Unmarshal(msg, &payload)
		if err != nil {
			log.Println(" Error decoding message:", err)
			continue
		}

		// Retrieve conversation details from the database
		var conversation models.DirectConversation
		if err := database.DB.First(&conversation, "id = ?", payload.DirectConversationID).Error; err != nil {
			log.Println(" Conversation not found:", err)
			continue
		}

		// Determine the recipient based on `SenderID`
		var receiverID uuid.UUID
		if payload.SenderID == conversation.User1ID {
			receiverID = conversation.User2ID
		} else {
			receiverID = conversation.User1ID
		}

		// Create a new message
		var directConversationID *uuid.UUID
		if payload.DirectConversationID != nil {
			directConversationID = payload.DirectConversationID
		}

		message := models.Message{
			ID:                   uuid.New(),
			DirectConversationID: directConversationID,
			SenderID:             payload.SenderID,
			Content:              payload.Content,
			CreatedAt:            time.Now(),
		}

		// Send the message to the recipient if they are online
		mu.Lock()
		receiverConn, receiverOnline := clients[receiverID]
		mu.Unlock()

		if receiverOnline {
			err := receiverConn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println(" Error sending message:", err)
				receiverConn.Close()
				mu.Lock()
				delete(clients, receiverID)
				mu.Unlock()
			}
		} else {
			log.Println("📌 User", receiverID, "is offline. Message will be delivered later.")
		}

		// Send message to the processing channel after sending it
		messageChannel <- message
	}
}
