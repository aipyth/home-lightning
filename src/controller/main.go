package main

import (
	"fmt"
	"sync"

	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"
)

const sleepTime = 50
const PLACE = "ceil"

type State struct {
	sync.Mutex
	Mode       string
	Brightness int
	Color      uint32
}

func serveLedUpdate(l *Leds, st *State, opt *ws2811.Option) {
	const sleepTime = 20
	modesTrigers := map[string]func(){
		"oneline":         onelineInitializer(l, st),
		"off":             offInitializer(l, st),
		"meteor":          meteorInitializer(l, st),
		"drops":           dropsInitializer(l, st),
		"fire":            fireInitializer(l, st),
		"2 soft colors":   twoSoftColorsInitializer(l, st),
		"oneline rainbow": onelineRainbowInitializer(l, st),
	}
	// clear all strips
	modesTrigers["off"]()

	for {
		l.SetBrightness(0, st.Brightness)
		l.SetBrightness(1, st.Brightness)

		if modesTrigers[st.Mode] != nil {
			modesTrigers[st.Mode]()
		}

		if err := l.Render(); err != nil {
			fmt.Println("ERROR", err)
		}
		if err := l.Wait(); err != nil {
			fmt.Println("ERROR", err)
		}
		// time.Sleep(time.Millisecond * sleepTime)
	}
}

func main() {
	fmt.Println("INFO Started")

	state := &State{
		Mode:       "oneline",
		Brightness: 255,
		Color:      0x000000,
	}

	go startRetrievingServerInfo(state)

	opt := ws2811.DefaultOptions
	opt.Channels = []ws2811.ChannelOption{
		{
			GpioPin:  18,
			LedCount: 246,
		},
		{
			GpioPin:  13,
			LedCount: 252,
		},
	}

	ws, err := ws2811.MakeWS2811(&opt)
	if err != nil {
		fmt.Println("ERROR", err)
		panic(err)
	}
	err = ws.Init()
	if err != nil {
		fmt.Println("ERROR", err)
		panic(err)
	}
	defer ws.Fini()

	leds := &Leds{ws}
	serveLedUpdate(leds, state, &opt)
}
