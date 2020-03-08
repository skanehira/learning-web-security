package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func printLog(msg string) {
	log.Output(2, msg)
}

func init() {
	log.SetFlags(log.Lshortfile)
	var err error
	db, err = sql.Open("sqlite3", "test.db")
	if err != nil {
		printLog(err.Error())
	}

	todoDDL := "create table if not exists todos (id integer primary key autoincrement, todo varchar(255))"
	if _, err := db.Exec(todoDDL); err != nil {
		printLog(err.Error())
	}

	userDDL := "create table if not exists users (name varchar(255), password varchar(255))"
	if _, err := db.Exec(userDDL); err != nil {
		printLog(err.Error())
	}

	addTestData()
}

func addTestData() {
	stmt, err := db.Prepare("insert into users (name, password) values (?, ?)")
	if err != nil {
		printLog(err.Error())
	}

	testData := map[string]string{
		"gorilla": "12344",
		"cat":     "!5698709",
	}

	for name, pass := range testData {
		if _, err := stmt.Exec(name, pass); err != nil {
			printLog(err.Error())
		}
	}
}

type Todo struct {
	ID   string
	Todo string
}

func createError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		sql := "select * from todos where id = " + r.URL.Query().Get("id")
		printLog(sql)

		rows, err := db.Query(sql)
		if err != nil {
			printLog(err.Error())
			createError(w, err)
			return
		}
		defer rows.Close()

		todos := []Todo{}
		for rows.Next() {
			var (
				id   string
				todo string
			)

			if err := rows.Scan(&id, &todo); err != nil {
				printLog(err.Error())
				createError(w, err)
				return
			}

			todos = append(todos, Todo{ID: id, Todo: todo})
		}

		if err := json.NewEncoder(w).Encode(&todos); err != nil {
			printLog(err.Error())
			createError(w, err)
			return
		}
	})

	log.Println("start http server :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
