package cli

import (
	"github.com/golang-migrate/migrate/v4"
	"log"
)

func Up(){
	m, err := migrate.New(
		"file://../migrations/",
		"postgresql://postgres:9406715@34.124.218.97:5432/personal?sslmode=disable")
	log.Printf("%v %v", err,m)
	err =m.Up()
	log.Printf("%v", err)
}