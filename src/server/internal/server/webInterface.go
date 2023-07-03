package server

import (
	"log"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"server/internal/storage"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


type WebClient struct {
	conn	*websocket.Conn
	storage	*storage.Storage
	toWrite	chan []byte
}




func (c *WebClient) readPump() {

	defer c.conn.Close()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })


	var	start time.Time
	var duration time.Duration
	var message []byte
	var err error

	for {
		_, message, err = c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[websocket error] %v", err)
			}
			break
		}

		start = time.Now()

		WebInterfaceHandler(c, message)

		duration = time.Since(start)
		log.Printf("[websocket] %s - %s",
			buildColor(boldBright, bgWhite, red) + "|" +
				duration.String() + "|" + buildColor(reset),
			string(message))
	}
}

func (c *WebClient) writePump() {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg := <-c.toWrite:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("[websocket write]", err)
			}

		case <-pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[websocket ping] %v", err)
				return
			}
		}
	}
}

func (s *server) serveWebInterface(ctx *fasthttp.RequestCtx) {
	err := upgrader.Upgrade(ctx, func (conn *websocket.Conn) {
		log.Println("[websocket] Opened connection to", conn.RemoteAddr())
		client := &WebClient{
			conn:		conn,
			storage:	s.Storage,
			toWrite:	make(chan []byte),
		}

		go client.readPump()
		client.writePump()
	})
	if err != nil {
		log.Printf("[websocket upgrader]", err)
	}
}
