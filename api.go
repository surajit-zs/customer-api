package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Customer struct {
	ID      int
	Name    string
	Phone   int
	Address string
}

func main() {
	var err error
	db, err = sql.Open("mysql", "surajit:Spore@0020@tcp(127.0.0.1:3306)/customers")
	if err != nil {
		log.Println(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Println(err)

	}

	http.HandleFunc("/customer", handler)
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		get(w, r)
	case http.MethodPost:
		post(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
}

// post reads the JSON body and inserts in the database
func post(w http.ResponseWriter, r *http.Request) {
	var customer Customer

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	err = json.Unmarshal(body, &customer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	_, err = db.Exec("INSERT INTO Customer (ID ,Name ,Phone Address) VALUES(?,?,?,?)", customer.ID, customer.Name, customer.Phone, customer.Address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	w.WriteHeader(http.StatusCreated)
}

// get retrieves the data from database and writes data as a JSON.
func get(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * from Customer;")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	defer func() {
		err := rows.Err()
		if err != nil {
			log.Println(err)
		}
	}()

	defer rows.Close()

	customers := make([]Customer, 0)

	for rows.Next() {
		var customer Customer

		err = rows.Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.Address)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		customers = append(customers, customer)

	}

	resp, err := json.Marshal(customers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}
