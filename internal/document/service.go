package document

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type DocumentService struct {
	upgrader websocket.Upgrader
}

func NewDocumentService() *DocumentService {
	return &DocumentService{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Allow all origins (adjust for production)
				return true
			},
		},
	}
}

func (d *DocumentService) UpgradeConnection(ctx *gin.Context) {
	conn, err := d.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket connection established")

	// Listen for messages from the client
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Echo message back to the client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}
