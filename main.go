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
	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/chaussette", handler.WSHandler)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	go handler.DispatchRoomMessage()
	go handler.DispatchPingMessage()

	fmt.Println(http.ListenAndServe(":"+port, router))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

}
