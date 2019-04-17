package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type (
	Message struct {
		id      int    `json:"id,omitempty"`
		done    bool   `json:"done"`
		message string `json:"message,omitempty"`
	}
	Messages      []Message
	ClientRequest struct {
		clientId int `json:"client_id,omitempty"`
		Message  `json:"message,omitempty"`
	}
	ClientResponse struct {
		Messages `json:"messages,omitempty"`
	}
)

var upgrader websocket.Upgrader

var messages Messages
var messageID int

func main() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":4001", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		clientsReq := &ClientRequest{}
		err := conn.ReadJSON(clientsReq)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("message from client: ", clientsReq)
		clientResp := &ClientResponse{}
		messageID++
		// message id
		clientsReq.Message.id = messageID
		// append to an array
		messages = append(messages, clientsReq.Message)
		// updated list
		clientResp.Messages = messages
		conn.WriteJSON(clientResp)
	}
}
