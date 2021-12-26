package rest

import (
	"dongpham/model"
	"dongpham/utils"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Api struct {
}

func (api *Api) BuildFuncApi(handler func(request *model.ApiRequest) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// ...
		// do something
		// ...
		request := &model.ApiRequest{
			Body: req.Body,
		}
		res, err := handler(request)
		if err != nil {
			logrus.Errorf("error when handling request %v %v", req.URL.Path, err)
			utils.ResponseResultAPIError(
				&model.ApiResponse{
					Code: 999,
					Data: "InternalError",
				}, w)
			return
		}
		utils.ResponseResultAPIError(
			&model.ApiResponse{
				Code: 0,
				Data: res,
			}, w)
	}
}

func RegisterRoutes(router *mux.Router) *mux.Router {
	//router.Methods().
	////router.Methods("POST").Path("/v1/content_json").HandlerFunc(GetTest)
	router = RegisterUserApi(router)
	router = RegisterGalleryApi(router)
	router = RegisterPostApi(router)

	return router
}
