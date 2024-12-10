package main

import (
	"database/sql"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Topic struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Created     string `json:"created"`
	Author      string `json:"author"`
	Profile     string `json:"profile"`
}

func main(){
	db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/test")
	if err != nil {
		fmt.Println("DB 열기 실패: ",err)
	}
	fmt.Println("Connected to MySQL database!")
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		rows, err := db.Query("SELECT id, title, description, created, author, profile FROM topic")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var data []Topic
        	for rows.Next() {
			var d Topic
                	if err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.Created, &d.Author, &d.Profile); err != nil {
                        	http.Error(w, err.Error(), http.StatusInternalServerError)
                        	return
                	}
			data = append(data,d)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(data); err !=nil {
			http.Error(w,err.Error(), http.StatusInternalServerError)
		}
        })
	fmt.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
