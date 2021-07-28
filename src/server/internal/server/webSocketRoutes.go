package server

import (
	"log"
	"bytes"
)

type ClientRequest [][]byte

var websocketRoutes = map[string]func(*WebClient, ClientRequest) {
	"schema": func(c *WebClient, r ClientRequest)  {
		msg := []byte(`create-mode;mode-name
get-modes
remove-mode;mode-name
create-place;place-name
get-places
remove-place;place-name
update-place;place-name;mode;color;brightness`)
		c.toWrite <- msg
	},

	// TODO some handler's logic repeats so make a fucking mixin whatever
	"create-mode": func(c *WebClient, r ClientRequest) {
		if len(r) <= 1 {
			c.toWrite <- []byte("mode name not provided")
			return
		}
		err := c.storage.AddMode(string(r[1]))
		if err != nil {
			log.Println("[websocket create-mode]", err)
			c.toWrite <- []byte("!server error")
			return
		}
		c.toWrite <- []byte("!mode created")
	},
	"get-modes": func (c *WebClient, r ClientRequest) {
		modesStr := c.storage.GetModes()
		modes := make([][]byte, len(modesStr) + 1)
		modes[0] = []byte("modes")
		for i := 0; i < len(modesStr); i++ {
			modes[i+1] = []byte(modesStr[i])
		}

		c.toWrite <- bytes.Join(modes, []byte(";"))
	},
	"remove-mode": func (c *WebClient, r ClientRequest) {
		if len(r) <= 1 {
			c.toWrite <- []byte("mode name not provided")
			return
		}
		err := c.storage.RemoveMode(string(r[1]))
		if err != nil {
			log.Println("[websocket remove-mode]", err)
			c.toWrite <- []byte("!server error")
			return
		}
		c.toWrite <- []byte("!mode removed")
	},

	"create-place": func (c *WebClient, r ClientRequest) {
		if len(r) <= 1 {
			c.toWrite <- []byte("place name not provided")
			return
		}
		err := c.storage.AddPlace(string(r[1]))
		if err != nil {
			log.Println("[websocket create-place]", err)
			c.toWrite <- []byte("!server error")
			return
		}
		c.toWrite <- []byte("!place created")
	},
	"get-places": func(c *WebClient, r ClientRequest) {
		placesStr := c.storage.GetPlaces()
		places := make([][]byte, len(placesStr) + 1)
		places[0] = []byte("places")
		for i := 0; i < len(placesStr); i++ {
			places[i+1] = []byte(placesStr[i])
		}

		c.toWrite <- bytes.Join(places, []byte(";"))
	},
	"remove-place": func(c *WebClient, r ClientRequest) {
		if len(r) <= 1 {
			c.toWrite <- []byte("!place name not provided")
			return
		}
		err := c.storage.RemovePlace(string(r[1]))
		if err != nil {
			log.Println("[websocket remove-place]", err)
			c.toWrite <- []byte("!server error")
			return
		}
		c.toWrite <- []byte("!place removed")
	},
	"update-place": func(c *WebClient, r ClientRequest) {
		if len(r) < 3 {
				c.toWrite <- []byte("!place not provided")
				return
		}
		place := string(r[1])
		args := r[2:]
		argsStr := make([]string, len(args))
		for i, v := range args { argsStr[i] = string(v) }
		err := c.storage.UpdatePlace(place, argsStr)
		if err != nil {
			log.Println("[websocket update-place]", err)
			c.toWrite <- []byte("!server error")
			return
		}
	},
	"get-place": func (c *WebClient, r ClientRequest) {
		if len(r) < 2 {
			c.toWrite <- []byte("!place not provided")
			return
		}
		place := string(r[1])
		valuesStr := c.storage.GetPlace(place)
		values := make([][]byte, len(valuesStr) + 1)
		values[0] = r[1]
		for i, v := range valuesStr { values[i+1] = []byte(v) }
		msg := bytes.Join(values, []byte(";"))
		c.toWrite <- msg
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
