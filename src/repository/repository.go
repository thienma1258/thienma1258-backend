package repository

import (
	"database/sql"
	"dongpham/config"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var dbConn *sql.DB
var GalleryRepo *GalleryRepository
var PostRepo *PostRepository

func init() {
	var err error
	dbConn, err = sql.Open("postgres", config.DBConnection)
	if err != nil {
		log.Errorln(err)
	}
	initRepository()
}

func initRepository() {
	GalleryRepo, _ = NewGalleryRepository(dbConn)
	PostRepo, _ = NewPostRepository(dbConn)

}

func CloseConn() {
	err := dbConn.Close()
	if err != nil {
		log.Errorln(err)
	}
}
