package main

import (
	"fmt"
	"os"
)

func main() {
	addr := os.Getenv("BACKEND_ADDR")
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	pswd := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s", host, user, pswd, db)
	var server Server
	err := server.Run(dsn, addr)
	if err != nil {
		panic(err)
	}
}
