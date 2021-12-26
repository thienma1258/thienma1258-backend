package rest

import (
	"dongpham/services"
	"github.com/gorilla/mux"
	"net/http"
)

type UserAPI struct {
	UserServices *services.UserServices
}

func (userApi *UserAPI) GetAllUser(w http.ResponseWriter, r *http.Request) {
	return
}

func RegisterUserApi(router *mux.Router) *mux.Router {
	User := UserAPI{}
	router.Methods("GET").Path("v0/Users").HandlerFunc(User.GetAllUser)
	return router
}
