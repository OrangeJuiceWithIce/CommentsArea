package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Comment struct {
	ID             int       `json:"id"`
	UserName       string    `json:"name"`
	CommentTime    time.Time `json:"time"`
	CommentContent string    `json:"content"`
}

func main() {
	r := http.NewServeMux()

	r.HandleFunc("/comment/add", HandleAddComment)
	r.HandleFunc("/comment/get", HandleGetComments)
	r.HandleFunc("/comment/delete", HandleDeleteComment)

	http.ListenAndServe(":8080", r)
}

func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var NewComment Comment
	err := json.NewDecoder(r.Body).Decode(&NewComment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("Comment", NewComment)
	defer r.Body.Close()

	db, err := sql.Open("mysql", "root:liuhan73@tcp(localhost:3306)/myfirstbase")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = InsertComment(db, NewComment.UserName, time.Now(), NewComment.CommentContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(NewComment)
}

func HandleGetComments(w http.ResponseWriter, r *http.Request) {

}

func HandleDeleteComment(w http.ResponseWriter, r *http.Request) {

}

func InsertComment(db *sql.DB, UserName string, CommentTime time.Time, CommentContent string) error {
	_, err := db.Exec("insert into comments(name,time,content) values(?,?,?)", UserName, CommentTime, CommentContent)
	if err != nil {
		return err
	}
	return nil
}
