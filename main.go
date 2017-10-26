package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/xwjdsh/2048-ai/ai"
	"github.com/xwjdsh/2048-ai/grid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var logo = `
██████╗  ██████╗ ██╗  ██╗ █████╗        █████╗ ██╗
╚════██╗██╔═████╗██║  ██║██╔══██╗      ██╔══██╗██║
 █████╔╝██║██╔██║███████║╚█████╔╝█████╗███████║██║
██╔═══╝ ████╔╝██║╚════██║██╔══██╗╚════╝██╔══██║██║
███████╗╚██████╔╝     ██║╚█████╔╝      ██║  ██║██║
╚══════╝ ╚═════╝      ╚═╝ ╚════╝       ╚═╝  ╚═╝╚═╝
`

var (
	addr   = flag.String("addr", ":8080", "http service address")
	online int32
)

func main() {
	fmt.Println(logo)
	flag.Parse()
	static := http.FileServer(http.Dir("./2048"))
	http.Handle("/js/", static)
	http.Handle("/style/", static)
	http.Handle("/meta/", static)
	http.Handle("/favicon.ico", static)

	indexTpl := template.Must(template.ParseFiles("./2048/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		indexTpl.Execute(w, nil)
	})
	http.HandleFunc("/compute", compute)

	log.Printf("Service started on \x1b[32;1m%s\x1b[32;1m\x1b[0m\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func compute(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err.Error())
		return
	}
	defer func() {
		if conn != nil {
			conn.Close()
		}
		online = atomic.AddInt32(&online, -1)
		log.Println(online)
	}()
	online = atomic.AddInt32(&online, 1)
	log.Println(online)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			break
		}
		g := &grid.Grid{}
		if err = json.Unmarshal(p, g); err != nil {
			break
		}
		a := &ai.AI{Grid: g}
		dire := a.Search()
		result := map[string]grid.Direction{"dire": dire}
		p, _ = json.Marshal(result)
		if err := conn.WriteMessage(messageType, p); err != nil {
			break
		}
	}
}
