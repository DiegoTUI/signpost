// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"gopkg.in/mgo.v2/bson"

	"path"

	"github.com/DiegoTUI/signpost/db"
	"github.com/DiegoTUI/signpost/models"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var (
	addr = flag.String("addr", "127.0.0.1:8080", "http service address")
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

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
		query := bson.M{"name": message}
		cities := []models.City{}

		err := db.Find(query, models.City{}.Collection(), &cities)
		if err != nil {
			log.Println("error in find", err)
			continue
		}

		if err := ws.WriteJSON(cities); err != nil {
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
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

var upgrader = websocket.Upgrader{}

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	reader := make(chan string, 100)

	stdoutDone := make(chan struct{})
	go pumpStdout(ws, reader, stdoutDone)
	go ping(ws, stdoutDone)

	pumpStdin(ws, reader)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path, r.Method)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	_, filename, _, _ := runtime.Caller(0)
	http.ServeFile(w, r, path.Dir(filename)+"/home.html")
}

func main() {
	// read environment
	var env string
	flag.StringVar(&env, "env", env, "Environment: 'development' or 'production'")
	flag.Parse()

	if env != "production" {
		env = "development"
	}

	// read config
	viper.SetConfigName("app")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found...")
		os.Exit(1)
	}

	dbhost := viper.GetString(env + ".dbhost")
	dbname := viper.GetString(env + ".dbname")

	// connect to the DB
	log.Println("Connecting to mongo")
	err = db.Connect(dbhost, dbname)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	// handle requests
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", serveWs)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
