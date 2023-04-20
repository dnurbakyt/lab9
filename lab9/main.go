package main

import (
	"log"

	"database/sql"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	db *sql.DB
}

func database() server {
	database, _ := sql.Open("sqlite3", "carsharing.db")
	server := server{db: database}
	return server
}

func (s *server) rent(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		fn := r.FormValue("fn")
		ln := r.FormValue("ln")
		cm := r.FormValue("cm")
		price := r.FormValue("price")
		hours := r.FormValue("hours")

		_, err := s.db.Exec("INSERT INTO carsharing(firstName, lastName, carModel, price, hours) VALUES ($1, $2, $3, $4, $5)", fn, ln, cm, price, hours)

		if err != nil {
			log.Fatal(err)
		}

		cars := map[string]interface{}{"carModel": cm, "text": "was rented!"}
		tmpl, _ := template.ParseFiles("static/rent.html")
		tmpl.Execute(w, cars)
		return
	}

	t, _ := template.ParseFiles("static/rent.html")
	t.Execute(w, nil)

}

func main() {
	s := database()

	fileServer := http.FileServer(http.Dir("./static"))

	http.Handle("/", fileServer)
	http.HandleFunc("/rent", s.rent)

	defer s.db.Close()
	http.ListenAndServe(":8080", nil)

}
