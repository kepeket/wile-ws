package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/wile-ws/handler"
)

func main() {
	log.SetOutput(os.Stderr)
	router := mux.NewRouter()
	router.HandleFunc("/chaussette", handler.WSHandler)

	go handler.DispatchRoomMessage()
	go handler.DispatchPingMessage()

	fmt.Println(http.ListenAndServe(":8844", router))
}
