package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Cant connect to websocket: ", err.Error())
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Can't read message: ", err.Error())
			return
		}

		response := reverse(string(p))

		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			log.Fatal("Error writing message: ", err.Error())
			return
		}
	}

}

func main() {
	fmt.Println("Starting server...")

	http.HandleFunc("/echo", echoHandler)
	http.Handle("/",
		http.FileServer(http.Dir("./template")))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error: " + err.Error())
	}
}
