package main

import (
	"ecommerce/auth"
	"ecommerce/database"
	"ecommerce/utility"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello!")
}

func main() {
	http.HandleFunc("/", helloHandler)
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("error when load env file with error", err.Error())
	}
	db, err := database.ConnectPostgres(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		panic(err)
	}
	if db == nil {
		panic("db not conncected")
	}
	// ======= UPDATE =======
	utility.InitToken("INI_SECRET", 1*60)
	// ======================
	router := chi.NewRouter()

	auth.Register2(router, db)

	// ======= UPDATE =======
	auth.Register2(router, db)
	// ======================
	const port = ":8000"
	log.Println("Server running at port", port)
	http.ListenAndServe(port, router)
}
