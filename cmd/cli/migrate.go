package cli

import (
	"github.com/golang-migrate/migrate/v4"
	"log"
)

func Up(){
	m, err := migrate.New(
		"file://../migrations/",
		"?")
	log.Printf("%v %v", err,m)
	err =m.Up()
	log.Printf("%v", err)
}