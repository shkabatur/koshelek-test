package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

type fromTo struct {
	From time.Time `json:"timeFrom"`
	To   time.Time `json:"timeTo"`
}

func wsConnect(w http.ResponseWriter, r *http.Request) {
	var err error
	wsConnection, err = upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("[wsConnect]: upgrade:", err)
		return
	}

	log.Println("[wsConnect]: new clinet connect over web sockets: ", r.RemoteAddr)
}

func home(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("templates/home.html")
	tpl.Execute(w, nil)
}

func client2(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("templates/client2.html")
	tpl.Execute(w, "ws://"+r.Host+"/ws")
}

func client3(w http.ResponseWriter, r *http.Request) {
	tpl, _ := template.ParseFiles("templates/client3.html")
	tpl.Execute(w, nil)
}

func receiveMessage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		log.Println("[receiveMessage]: Method ", r.Method, " is not allowed!")
		w.WriteHeader(405) // Method Not Allowed.
		return
	}

	// Read request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[receiveMessage]: Body read error, %v", err)
		w.WriteHeader(500) // Internal Server Error.
		return
	}

	m := message{}
	// Parse message
	if err = json.Unmarshal(body, &m); err != nil {
		log.Printf("[receiveMessage]: Body parse error, %v", err)
		w.WriteHeader(400) // Bad Request.
		return
	}
	m.Time = time.Now()
	w.WriteHeader(200)
	log.Println("[receiveMessage]: You have one new message: ", m)

	err = writeMessageToDB(m)
	if err != nil {
		log.Println("[receiveMessage]: Error while writing message to db: ", err)
		return
	}

	if wsConnection != nil {
		js, err := json.Marshal(&m)
		if err != nil {
			log.Println("[receiveMessage]: Cant marshal message in receiveMessage function: ", err)
			return
		}
		err = wsConnection.WriteJSON(string(js))
		if err != nil {
			log.Println("[receiveMessage]: Can't send message over web socket: ", err)
			return
		}
	}

}

func getMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("[getMessages]: Method ", r.Method, " is not allowed!")
		w.WriteHeader(405) // Method Not Allowed.
		return
	}

	// Read request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[getMessages]: Body read error, %v", err)
		w.WriteHeader(500) // Internal Server Error.
		return
	}

	ft := fromTo{}
	// Parse message
	if err = json.Unmarshal(body, &ft); err != nil {
		log.Printf("[getMessages]: Body parse error, %v", err)
		w.WriteHeader(400) // Bad Request.
		return
	}

	messages, err := getMessagesFromDB(ft.From, ft.To)
	if err != nil {
		return
	}

	js, err := json.Marshal(messages)
	if err != nil {
		fmt.Println("[getMessages]: Error while joson.Marshal(messages): ", err)
	}

	w.Header().Set("[getMessages]: Content-Type", "application/json")
	w.Write(js)
}
