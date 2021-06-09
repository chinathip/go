package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Actor struct {
	Address_id int    `json:"address_id"`
	Address    string `json:"address"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "chinathip"
	password = "root"
	dbname   = "test_db"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT address_id,address FROM address ORDER BY address_id DESC")
	if err != nil {
		log.Fatal(err)
	}

	var people []Actor

	for rows.Next() {
		var actor Actor
		rows.Scan(&actor.Address_id, &actor.Address)
		people = append(people, actor)
	}

	peopleBytes, _ := json.MarshalIndent(people, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(peopleBytes)

	defer rows.Close()
	defer db.Close()
}

func main() {
	http.HandleFunc("/", GETHandler)
	log.Fatal(http.ListenAndServe(":2500", nil))

}
