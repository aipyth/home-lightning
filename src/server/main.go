package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	tb "gopkg.in/tucnak/telebot.v2"
)

//  ________________________
// |------------------------|
// |                       -|
// |                       -|___________________________________
// |                       ------------------------------------ |
// |                                                            |
// |                                                            |
// |                                                            |
// |                                                            |
// |------------------------                                    |
// |_______________________-                                   -|
// |-                     |-                                   -|
// |-                     |-                                   -|
// |-                     |-                                   -|
// |-                     |                                    -|
// |-                     |                                    -|
// |-                     |                      -------------- |
// |____________________________________________________________|

type LEDLine struct {
	Mode       string
	Place      string
	Color      string
	Brightness float32
}

var Modes = []string{
	"meteor",
	"oneline",
	"off",
}

func StartBot(leds []LEDLine) {
	b, err := tb.NewBot(tb.Settings{
		Token:  "1467814308:AAES_W2E8Nudwd_tcHk2nwjAJzVERWctZUM",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/leds", func(m *tb.Message) {

		ledsKeyboard := make([][]tb.InlineButton, 0)

		for _, l := range leds {
			ledsKeyboard = append(ledsKeyboard, []tb.InlineButton{
				tb.InlineButton{
					Unique: l.Place,
					Text:   "Line " + l.Place,
				},
			})
		}

		b.Send(m.Sender, "Hello World!", &tb.ReplyMarkup{
			InlineKeyboard: ledsKeyboard,
		})
	})

	b.Handle(tb.OnCallback, func(c *tb.Callback) {
		//ledActions := []string{
		//	"brightness",
		//	"mode",
		//}
		//
		//action := "none"
		//for _, l := range leds {
		//	if l.Place == c.Data {
		//
		//	}
		//}
		//
		//
		//if action == "place" {
		//	b.Send(c.Sender, "What to do?", &tb.ReplyMarkup{
		//
		//	})
		//} else if action == "led" {
		//
		//}

		b.Respond(c)
	})

	b.Start()
}

func StartServer(leds []LEDLine) {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/modes", func(c *gin.Context) {
		c.JSON(200, Modes)
	})

	r.GET("/leds", func(c *gin.Context) {
		c.JSON(200, leds)
	})
	r.POST("/leds", func(c *gin.Context) {
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("ERROR", err)
		}
		err = json.Unmarshal(jsonData, &leds)
		if err != nil {
			log.Println("ERROR", err)
		}
	})

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"leds": leds,
		})
	})
	r.Run()
}

func main() {
	LEDs := []LEDLine{
		LEDLine{
			Mode:       "oneline",
			Place:      "ceil",
			Color:      "#ffffff",
			Brightness: float32(0),
		},
		LEDLine{
			Mode:       "off",
			Place:      "kitchen",
			Color:      "#ffffff",
			Brightness: float32(0),
		},
	}

	go StartBot(LEDs)
	StartServer(LEDs)
}
