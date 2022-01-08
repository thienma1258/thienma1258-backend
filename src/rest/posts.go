package rest

import (
	"dongpham/model"
	"dongpham/repository"
	"dongpham/services"
	"github.com/gorilla/mux"
	"net/http"
)

type PostAPI struct {
	postService *services.PostServices
	Api
}

func (postAPi *PostAPI) GetIDs(request *model.ApiRequest) (interface{}, error) {
	ids, err := postAPi.postService.GetAllPostIDs()
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (postAPi *PostAPI) Create(request *model.ApiRequest) (interface{}, error) {
	var createRequest *model.Post
	err := json.NewDecoder(request.Body).Decode(&createRequest)
	if err != nil {
		return nil, err
	}

	err = postAPi.postService.Create(*createRequest)

	return nil, err
}

func (postAPi *PostAPI) Update(request *model.ApiRequest) (interface{}, error) {
	var createRequest *model.Post
	err := json.NewDecoder(request.Body).Decode(&createRequest)
	if err != nil {
		return nil, err
	}

	err = postAPi.postService.Update(*createRequest)

	return nil, err
}

func RegisterPostApi(router *mux.Router) *mux.Router {
	post := PostAPI{
		postService: services.NewPostServices(repository.PostRepo),
	}
	router.Methods("GET").Path("/v0/posts/ids").HandlerFunc(post.BuildFuncApi(post.GetIDs))
	router.Methods("POST").Path("/v0/posts").HandlerFunc(post.BuildFuncApi(post.Create))
	router.Methods(http.MethodPut).Path("/v0/posts/{id}").HandlerFunc(post.BuildFuncApi(post.Update))


	return router
}
