package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/joho/godotenv"
)

const port string = ":8080"

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}
func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		f(w, r)
	}
}

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func search(w http.ResponseWriter, req *http.Request) {

	apiKey := os.Getenv("GIANT_BOMB_API_KEY")
	query := req.URL.Query().Get("query")
	resp, err := http.Get("https://www.giantbomb.com/api/search/?api_key=" + apiKey + "&format=json&query=" + query + "&resources=game")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (env *Env) register(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := req.URL.Query().Get("username")
	email := req.URL.Query().Get("email")
	password := req.URL.Query().Get("password")

	addUser(env.db, username, email, password)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

type Env struct {
	db *sql.DB
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env variables")
	}

	driver := "mysql"
	connectionString := "mies:secret@(127.0.0.1:3306)/test?parseTime=true"

	db, err := sql.Open(driver, connectionString)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	env := &Env{db: db}

	{ // Create a new table
		query := `
		    CREATE TABLE users (
			id INT AUTO_INCREMENT,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME,
			PRIMARY KEY (id)
		    );`

		if _, err := db.Exec(query); err != nil {
			log.Fatal(err)
		}
	}
	http.HandleFunc("/hello", logging(hello))
	http.HandleFunc("/headers", logging(headers))
	http.HandleFunc("/search", logging(search))
	http.HandleFunc("/register", logging(env.register))
	fmt.Printf("running server on port %s \n", port)
	http.ListenAndServe(port, nil)
}
