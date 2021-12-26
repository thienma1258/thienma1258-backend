package rest

import (
	"dongpham/repository"
	"dongpham/services"
	"dongpham/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type Gallery struct {
	GalleryServices *services.GalleryServices
}

func (galleryApi *Gallery) GetAllGalleryImages(w http.ResponseWriter, r *http.Request) {
	articles, err := galleryApi.GalleryServices.GetAllGallery()

	if err != nil {
		utils.ResponseWithCodeAndData(utils.ERROR_UNKNOWN_ERROR, []byte(err.Error()), w)
		return
	}

	jsonData, err := json.Marshal(articles)

	if err != nil {
		utils.ResponseError(utils.ERROR_INVALID_REQUEST, w)
		return
	}

	utils.ResponseResultByte(jsonData, w)
}

func RegisterGalleryApi(router *mux.Router) *mux.Router {
	Gallery := Gallery{
		GalleryServices: services.NewGalleryServices(repository.GalleryRepo),
	}
	router.Path("/").HandlerFunc(Gallery.GetAllGalleryImages)


	router.Methods("GET").Path("/api/v0/Gallery").HandlerFunc(Gallery.GetAllGalleryImages)
	return router
}
