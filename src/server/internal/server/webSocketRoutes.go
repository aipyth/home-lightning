package server

import (
	"log"
	"bytes"
)

type ClientRequest [][]byte

var websocketRoutes = map[string]func(*WebClient, ClientRequest) {
	"test": func(c *WebClient, r ClientRequest)  {
		c.toWrite <- []byte("realy test")
	},

	"create-mode": func(c *WebClient, r ClientRequest) {
		if len(r) <= 1 {
			c.toWrite <- []byte("mode name not provided")
			return
		}
		err := c.storage.AddMode(string(r[1]))
		if err != nil {
			log.Println("[websocket create-mode]", err)
			c.toWrite <- []byte("server error")
			return
		}
		c.toWrite <- []byte("created")
	},
	"get-modes": func (c *WebClient, r ClientRequest) {
		modesStr := c.storage.GetModes()
		modes := make([][]byte, len(modesStr))
		for i := 0; i < len(modesStr); i++ {
			modes[i] = []byte(modesStr[i])
		}

		c.toWrite <- bytes.Join(modes, []byte(";"))
	},
}

func WebInterfaceHandler(c *WebClient, message []byte) {
	args := bytes.Split(message, []byte(";"))
	request := string(args[0])

	handler := websocketRoutes[request]
	if handler == nil {
		log.Printf("[websocket:%v] No found handler for %s", c.conn.RemoteAddr(), request)
		c.toWrite <- []byte("No found handler")
	} else {
		handler(c, args)
	}
}
