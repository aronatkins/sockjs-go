package main

import (
	"log"
	"net/http"

	"github.com/igm/pubsub"
	"github.com/igm/sockjs-go/sockjs"
)

var chat pubsub.Publisher

func main() {
	http.Handle("/echo/", sockjs.NewHandler("/echo", sockjs.DefaultOptions, echoHandler))
	http.Handle("/", http.FileServer(http.Dir("web/")))
	log.Println("Server started")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func echoHandler(conn sockjs.Session) {
	log.Println("new sockjs connection established")
	var closedSession = make(chan struct{})
	chat.Publish("[info] new participant joined chat")
	defer chat.Publish("[info] participant left chat")
	go func() {
		reader, _ := chat.SubChannel(nil)
		for {
			select {
			case <-closedSession:
				return
			case msg := <-reader:
				if err := conn.Send(msg.(string)); err != nil {
					return
				}
			}

		}
	}()
	for {
		if msg, err := conn.Recv(); err == nil {
			chat.Publish(msg)
			continue
		}
		break
	}
	close(closedSession)
	log.Println("sockjs connection closed")
}
