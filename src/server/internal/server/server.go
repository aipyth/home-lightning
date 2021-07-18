package server

import (
	"os"
	"sync"
	"log"
	"os/signal"
	"syscall"
	"time"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"server/internal/storage"
)
type server struct {
	HTTPServer	*fasthttp.Server
	Router		*router.Router
	Storage		*storage.Storage
}

func newServer() *server {
	router := router.New()
	httpServer := &fasthttp.Server{
		Handler:	router.Handler,
		ReadTimeout:          5 * time.Second,
		WriteTimeout:         10 * time.Second,
		MaxConnsPerIP:        500,
		MaxRequestsPerConn:   500,
		MaxKeepaliveDuration: 5 * time.Second,
	}
	return &server{
		HTTPServer:	httpServer,
		Router:		router,
	}
}


func (s *server) Start(address string) error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("hostname unavailable: %s", err)
	}

	log.Printf("%s - Web server starting on %v", hostname, address)
	log.Printf("%s - Press Ctrl+C to stop", hostname)


	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()

		err := s.HTTPServer.ListenAndServe(address)
		if err != nil {
			log.Println(err)
		}
		log.Println("Server goroutine finished")
	}()

	ticker := time.NewTicker(time.Second)

	for {
		select{
		case <-ticker.C:
			continue
		case <-osSignals:
			log.Println("Shutdown signal received")

			s.HTTPServer.Shutdown()

			wg.Wait()
			return nil
		}
	}
}


func Run() error {
	serv := newServer()
	serv.Routes()
	var err error
	serv.Storage, err = storage.NewStorage()
	if err != nil {
		log.Println("[storage]", err)
		return err
	}
	return serv.Start(":8080")
}
