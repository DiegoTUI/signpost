package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"

	"path"

	"github.com/DiegoTUI/signpost/db"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var (
	addr     = flag.String("addr", "127.0.0.1:8080", "http service address")
	upgrader = websocket.Upgrader{}
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	buffer := make(chan string, 100)

	stdoutDone := make(chan struct{})
	go pumpStdout(ws, buffer, stdoutDone)
	go ping(ws, stdoutDone)

	pumpStdin(ws, buffer)
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
