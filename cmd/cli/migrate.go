package cli

import (
	"github.com/golang-migrate/migrate/v4"
	"log"
)

func Up(){
	m, err := migrate.New(
		"file://../migrations/",
		"postgresql://postgres:9406715@localhost:5432/personalDB?sslmode=disable")
	log.Printf("%v %v", err,m)
	err =m.Up()
	log.Printf("%v", err)
}