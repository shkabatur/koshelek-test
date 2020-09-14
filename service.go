package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

const (
	host     = "postgres"
	port     = 5432
	user     = "postgres"
	password = "denis_denis"
	dbname   = "postgres"
)

type message struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
	SN   int       `json:"sn"`
}

var db *sql.DB
var wsConnection *websocket.Conn
var upgrader = websocket.Upgrader{}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("[main]: Can't open sql connection: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("[main]: Successfully connected to db!")

	err = createTable()
	if err != nil {
		log.Fatal("[main]:Can't create table: ", err)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/client2", client2)
	http.HandleFunc("/client3", client3)
	http.HandleFunc("/ws", wsConnect)
	http.HandleFunc("/messages", getMessages)
	http.HandleFunc("/send-message", receiveMessage)

	go sendMessageEvery(10 * time.Second)
	http.ListenAndServe("0.0.0.0:8080", nil)

}

// sendMessageEvery: send post request every time interval
func sendMessageEvery(interval time.Duration) {
	SN := 0
	tick := time.Tick(interval)
	for {
		SN++
		select {
		case <-tick:
			m := message{Text: "some message", SN: SN}

			rBody, err := json.Marshal(&m)
			if err != nil {
				log.Println("[sendMessageEvery]: Cant marshal message: ", err)
				return
			}

			resp, err := http.Post("http://localhost:8080/send-message", "application/json", bytes.NewBuffer(rBody))
			if err != nil {
				log.Println("[sendMessageEvery]: Error when try to send Post request: ", err)
				return
			}
			defer resp.Body.Close()
		}
	}

}
