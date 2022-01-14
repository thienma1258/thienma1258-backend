package rest

import (
	"dongpham/config"
	"dongpham/internal_errors"
	"dongpham/model"
	"dongpham/utils"
	"dongpham/version"
	"fmt"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Api struct {
}

func (api *Api) BuildFuncApi(handler func(request *model.ApiRequest) (interface{}, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		// ...
		// do something
		// ...
		query := mux.Vars(req)
		request := &model.ApiRequest{
			Body:     req.Body,
			Query:    query,
			URLQuery: req.URL.Query(),
		}
		res, err := handler(request)
		if err != nil {

			serr, ok := err.(*internal_errors.InternalError)
			if ok {
				utils.ResponseResultAPIError(
					&model.ApiResponse{
						Code: serr.Code,
						Data: serr.Message,
					}, w)
			} else {
				logrus.Errorf("error when handling request %v %v", req.URL.Path, err)
				utils.ResponseResultAPIError(
					&model.ApiResponse{
						Code: 999,
						Data: fmt.Sprintf("InternalError: %v", err),
					}, w)
				return
			}
		}
		utils.ResponseResultAPIError(
			&model.ApiResponse{
				Code: 0,
				Data: res,
			}, w)
	}
}

func RegisterRoutes(router *mux.Router) *mux.Router {
	router.Use(AppInfo(config.ServiceName, "INKR Global", version.Version))
	router.Methods(http.MethodOptions).HandlerFunc(CORS)
	////router.Methods("POST").Path("/v1/content_json").HandlerFunc(GetTest)
	router = RegisterUserApi(router)
	router = RegisterGalleryApi(router)
	router = RegisterPostApi(router)
	router = RegisterObjectMetaApi(router)
	return router
}

func CORS(writer http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	if strings.Index(origin, "phamdong.com") > 0 ||
		strings.Index(origin, "ngocdong.com") > 0 ||
		strings.Index(origin, "localhost") > 0 ||
		strings.ToUpper(r.Method) == "OPTIONS" {
		writer.Header().Set("Access-Control-Allow-Origin", origin)
		writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
		writer.Header().Set("Access-Control-Max-Age", "86400")
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, If-None-Match")
	}
	if strings.ToUpper(r.Method) == "OPTIONS" {
		writer.WriteHeader(204) // send the headers with a 204 response code.
	}
}

// AppInfo adds custom app-info to the response header
func AppInfo(app string, author string, version string) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			CORS(w, r)
			w.Header().Set("Author", author)
			w.Header().Set("App-Name", app)
			w.Header().Set("App-Version", version)
			if mhost := config.ServiceName; mhost != "" {
				w.Header().Set("Host", mhost)
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
