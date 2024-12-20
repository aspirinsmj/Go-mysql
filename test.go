package main

import (
	"database/sql"
	"net/http"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/jango", username, password, dbHost, dbPort)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("DB 열기 실패: ",err)
	}
	fmt.Println("Connected to MySQL database!")
	defer db.Close()

	err = db.Ping()
        if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}
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
