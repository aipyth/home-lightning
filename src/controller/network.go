package main

// #include <unistd.h>
// //#include <errno.h>
// //int usleep(useconds_t usec);
import "C"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
)

type ServerResponseNode struct {
	Mode       string
	Place      string
	Color      string
	Brightness float32
}

func toColor(c string) uint32 {
	r, _ := strconv.ParseInt(c[1:3], 16, 32)
	g, _ := strconv.ParseInt(c[3:5], 16, 32)
	b, _ := strconv.ParseInt(c[5:7], 16, 32)
	return uint32(r<<8 | g<<16 | b)
}

func startRetrievingServerInfo(st *State) {
	// frequency := time.Millisecond * 100
	// for {
	// 	<-time.Tick(frequency)
	// 	retrieveServerInfo(st)
	// }
	for {
		retrieveServerInfo(st)
		C.usleep(100000)
	}
}

func retrieveServerInfo(st *State) {
	const host = "192.168.1.7"

	resp, err := http.Get("http://" + host + "/leds")
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	serverResponse := make([]ServerResponseNode, 0)
	json.Unmarshal(body, &serverResponse)

	for _, v := range serverResponse {
		if v.Place == PLACE {
			st.Lock()
			st.Mode = v.Mode
			st.Brightness = int(math.Floor(float64(v.Brightness) * 255))
			st.Color = toColor(v.Color)
			st.Unlock()
		}
	}
}
