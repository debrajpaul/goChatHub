package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	socketio "github.com/googollee/go-socket.io"
)

type Client struct {
	ClientID int64 `json:"client_id,omitempty"`
}

var clients []Client

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		var clientInst Client
		fmt.Println("connected:", s.ID())
		// covert string to int 64
		n, err := strconv.ParseInt(s.ID(), 36, 64)
		if err == nil {
			fmt.Printf("%d of type %T", n, n)
		}
		clientInst.ClientID = n
		// client list
		clients = append(clients, clientInst)
		fmt.Println("The Clients are :=== ", clients)

		// for current client Id
		fmt.Println("userId:", clientInst)
		s.Emit("userId", clientInst)

		// for list of other client id
		var filterClient []Client
		for _, client := range clients {
			if client.ClientID != clientInst.ClientID {
				filterClient = append(filterClient, client)
			}
		}

		fmt.Println("userIdList:", filterClient)
		s.Emit("userIdList", filterClient)
		return nil
	})
	server.OnEvent("/", "notice", func(s socketio.Conn, listOfUser string, msg string) {
		fmt.Println("notice:", msg)
		//s.Emit("reply", "have "+msg)
		users := strings.Fields(listOfUser)
		s.BroadcastTo(users, "some:event", msg)
	})
	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})
	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, msg string) {
		fmt.Println("closed", msg)
	})
	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Serving at localhost:4001...")
	log.Fatal(http.ListenAndServe(":4001", nil))
}
