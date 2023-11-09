package main

import (
	"go-sample/api"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
)

var db *sqlx.DB

func init() {
	var err error

	gotenv.Load()
	db, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err)
	}

}
func main() {
	srv := api.NewServer(os.Getenv("PORT"), db)
	log.Println("Server running on port: ", os.Getenv("PORT"))
	log.Fatal(srv.Start())
	defer db.Close()
}
