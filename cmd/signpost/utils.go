package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/DiegoTUI/signpost/models"
	"github.com/gorilla/websocket"
)

var (
	pingErrors = make(map[string]uint8)
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Maximum message size allowed from peer.
	maxPingErrors = 3

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

func pumpStdin(ws *websocket.Conn, w chan string) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message)
		log.Println("message in", string(message))
		w <- string(message)
	}
}

func pumpStdout(ws *websocket.Conn, r chan string, done chan struct{}) {
	defer func() {
	}()

	for message := range r {
		log.Println("message out", message)
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		// go look in the db for a city
		query, err := models.NewSignpostQuery(message)
		if err != nil {
			ws.WriteJSON(bson.M{"error": err.Error()})
			continue
		}

		signpost, err := models.NewSignpost(query.Center, query.MinNumberOfSigns, query.MaxNumberOfSigns,
			query.MinDistance, query.MaxDistance, query.MinDifficulty, query.MaxDifficulty)
		if err != nil {
			ws.WriteJSON(bson.M{"error": err.Error()})
			continue
		}

		if err := ws.WriteJSON(signpost); err != nil {
			ws.Close()
			break
		}
	}

	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			wsID := ws.RemoteAddr().String()
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				_, ok := pingErrors[wsID]
				if ok {
					pingErrors[wsID]++
				} else {
					pingErrors[wsID] = 1
				}

				log.Println("ping:", err)

				if pingErrors[wsID] == maxPingErrors {
					log.Println("closing websocket:", wsID)
					close(done)
				}
			} else {
				delete(pingErrors, wsID)
			}
		case <-done:
			return
		}
	}
}
