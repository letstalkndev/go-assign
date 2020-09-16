package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Blog describe
type Blog struct {
	ID    int
	Title string
	Desc  string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "goassign"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:8889)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

//var tmpl = template.Must(template.ParseGlob("views/*"))

func index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM blog ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	blog := Blog{}
	blogs := []Blog{}

	for selDB.Next() {
		var id int
		var title, description string
		err = selDB.Scan(&id, &title, &description)
		if err != nil {
			panic(err.Error())
		}
		blog.ID = id
		blog.Title = title
		blog.Desc = description
		blogs = append(blogs, blog)
	}

	// tmpl.ExecuteTemplate(w, "index", res)
	js, err := json.Marshal(blogs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func insert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		title := r.FormValue("title")
		desc := r.FormValue("desc")
		insForm, err := db.Prepare("INSERT INTO blog(title, description) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(title, desc)
		log.Println("INSERT: title: " + title + " | desc: " + desc)
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/insert", insert)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
