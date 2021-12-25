package main

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	m, err := migrate.New(
		"file://../migrations/",
		"postgresql://postgres:9406715@localhost:5432/personalDB?sslmode=disable")
	log.Printf("%v %v", err,m)
	err =m.Up()
	log.Printf("%v", err)
}
