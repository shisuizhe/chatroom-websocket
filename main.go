package main

import (
	"chatroom/better_version"
	"chatroom/simple_version"
	"chatroom/util"
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func simple() error {
	simple_version.InitBroker()

	streamHandler, err := simple_version.NewMessageStreamHandler(&websocket.Upgrader{})
	if err != nil {
		return err
	}

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/streams/messages", streamHandler)

	return nil
}

func better() error {
	better_version.InitBroker(context.Background())

	streamHandler, err := better_version.NewMessageStreamHandler(&websocket.Upgrader{})
	if err != nil {
		return err
	}

	http.Handle("/", http.FileServer(http.Dir("public")))
	http.Handle("/streams/messages", streamHandler)

	return nil
}

func main() {
	// err := simple()
	err := better()

	if err != nil {
		log.Fatal(err)
	}

	util.LogInfo("starting server on port 8000")

	err = http.ListenAndServe(":8000", nil)
	if err != nil {
		util.LogErr("ListenAndServe", err)
	}
}
