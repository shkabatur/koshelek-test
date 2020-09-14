package main

import (
	"log"
	"time"
)

func createTable() error {
	createTableSQL := `
	create table if not exists messages (
		id				SERIAL PRIMARY KEY,
		message 		TEXT,
		t				TIMESTAMP,
		sn 				INTEGER
	)
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}
	return nil
}

// getMessages: get messages from db between timeFrom and timeTo
func getMessagesFromDB(timeFrom time.Time, timeTo time.Time) ([]message, error) {
	log.Println("[getMessagesFromDB]: Try to get messages between ", timeFrom, " and ", timeTo)
	messages := []message{}

	rows, err := db.Query("SELECT message, t, sn FROM messages WHERE t > $1 and t < $2", timeFrom.Local(), timeTo.Local())
	if err != nil {
		log.Println("[getMessagesFromDB]: Cant query messages: ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		m := message{}
		err = rows.Scan(&m.Text, &m.Time, &m.SN)
		if err != nil {
			log.Println("[getMessagesFromDB]: rows.Scan error: ", err)
			return nil, err
		}
		m.Time = m.Time.Local()
		messages = append(messages, m)
	}
	log.Println("[getMessagesFromDB]: Length of getted messages is :", len(messages))
	return messages, nil
}

// writeMessageToDB: take message, write it to database
func writeMessageToDB(m message) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert into messages ( message, t, sn ) values($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.Text, m.Time, m.SN)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
