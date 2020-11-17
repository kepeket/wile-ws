package handler

// // DispatchWorkoutMessage dispatch the room join message
// func DispatchWorkoutMessage() {
// 	for {
// 		val := <-roomBroadcast
// 		log.Println("new room message received")
// 		mess := fmt.Sprintf("someone has joined the room")
// 		// send to every client that is currently connected
// 		for client := range rooms[val.Name].Clients {
// 			err := client.WriteMessage(websocket.TextMessage, []byte(mess))
// 			if err != nil {
// 				log.Printf("Websocket error: %s", err)
// 				client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }
