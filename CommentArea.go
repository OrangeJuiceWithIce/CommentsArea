package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// 统一response结构
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// 评论结构
type Comment struct {
	ID             int       `json:"id"`
	UserName       string    `json:"name"`
	CommentTime    time.Time `json:"time"`
	CommentContent string    `json:"content"`
}

// 主函数，分配路由
func main() {
	r := http.NewServeMux()

	r.HandleFunc("/comment/add", ResponseHandler(HandleAddComment))
	r.HandleFunc("/comment/get", ResponseHandler(HandleGetComment))
	r.HandleFunc("/comment/delete", ResponseHandler(HandleDeleteComment))

	http.ListenAndServe(":8080", r)
}

// 把三个函数封装成一个函数，方便调用
type handler func(w http.ResponseWriter, r *http.Request) (data interface{}, err error)

// 添加评论,遇到错误直接返回
func HandleAddComment(w http.ResponseWriter, r *http.Request) (data interface{}, err error) {
	var NewComment Comment
	err = json.NewDecoder(r.Body).Decode(&NewComment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	fmt.Println("Comment", NewComment)
	defer r.Body.Close()

	db, err := sql.Open("mysql", "root:liuhan73@tcp(localhost:3306)/myfirstbase")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	err = InsertComment(db, NewComment.UserName, time.Now(), NewComment.CommentContent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return NewComment, nil
}

// 获取评论
func HandleGetComment(w http.ResponseWriter, r *http.Request) (data interface{}, err error) {
	db, err := sql.Open("mysql", "root:liuhan73@tcp(localhost:3306)/myfirstbase")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()
	queryvalue := r.URL.Query()
	page, _ := strconv.Atoi(queryvalue.Get("page"))
	size, _ := strconv.Atoi(queryvalue.Get("size"))

	type pagedata struct {
		Total    int       `json:"total"`
		Comments []Comment `json:"comments"`
	}
	var comments []Comment
	comments, err = GetComment(db, page, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	Pagedata := pagedata{len(comments), comments}
	return Pagedata, nil
}

// 删除评论
func HandleDeleteComment(w http.ResponseWriter, r *http.Request) (data interface{}, err error) {
	db, err := sql.Open("mysql", "root:liuhan73@tcp(localhost:3306)/myfirstbase")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()
	queryvalue := r.URL.Query()
	id, _ := strconv.Atoi(queryvalue.Get("id"))
	err = DeleteComment(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}
	return nil, nil
}

// 对handler函数进行封装，统一返回格式
func ResponseHandler(h handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		data, err := h(w, r)
		//遇到错误返回错误的response
		if err != nil {
			resp := Response{Code: 1, Msg: err.Error(), Data: nil}
			json.NewEncoder(w).Encode(resp)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		//否则返回成功的response
		resp := Response{Code: 0, Msg: "success", Data: data}
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(http.StatusOK)
	}
}

func InsertComment(db *sql.DB, UserName string, CommentTime time.Time, CommentContent string) error {
	_, err := db.Exec("insert into comments(name,time,content) values(?,?,?)", UserName, CommentTime, CommentContent)
	if err != nil {
		return err
	}
	return nil
}

func GetComment(db *sql.DB, page int, size int) ([]Comment, error) {
	var rows *sql.Rows
	var err error
	fmt.Println("page", page, "size", size)
	if size == -1 {
		rows, err = db.Query("select id,name,content from comments")
	} else {
		rows, err = db.Query("select id,name,content from comments limit ?,?", (page-1)*size, size)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		err = rows.Scan(&c.ID, &c.UserName, &c.CommentContent)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	fmt.Println(comments)
	return comments, nil
}

func DeleteComment(db *sql.DB, id int) error {
	fmt.Println("delete id", id)
	_, err := db.Exec("delete from comments where id=?", id)
	if err != nil {
		return err
	}
	return nil
}
