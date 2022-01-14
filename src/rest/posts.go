package rest

import (
	"dongpham/model"
	"dongpham/repository"
	"dongpham/services"
	"dongpham/utils"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type PostAPI struct {
	postService *services.PostServices
	Api
}

const POST_ID = "id"

func (postAPi *PostAPI) GetIDs(request *model.ApiRequest) (interface{}, error) {
	var publishedQuery *bool
	if len(request.URLQuery.Get("published")) > 0 {
		published, err := strconv.ParseBool(request.URLQuery.Get("published"))
		if err != nil {
			return nil, err
		}
		publishedQuery = &published
	}
	var orderDESC *bool
	if len(request.URLQuery.Get("orderDESC")) > 0 {
		rawOrderDESC, err := strconv.ParseBool(request.URLQuery.Get("orderDESC"))
		if err != nil {
			return nil, err
		}
		publishedQuery = &rawOrderDESC
	}
	ids, err := postAPi.postService.GetAllPostIDs(publishedQuery,orderDESC)
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

	result, err := postAPi.postService.Create(*createRequest)

	return result, err
}

func (postAPi *PostAPI) Update(request *model.ApiRequest) (interface{}, error) {
	var updateRequest *model.Post
	err := json.NewDecoder(request.Body).Decode(&updateRequest)
	if err != nil {
		return nil, err
	}
	rawID := request.Query[POST_ID]
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return nil, err
	}
	updateRequest.ID = utils.Int(id)
	err = postAPi.postService.Update(*updateRequest)

	return nil, err
}

func (postAPi *PostAPI) Delete(request *model.ApiRequest) (interface{}, error) {

	rawID := request.Query[POST_ID]
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return nil, err
	}
	err = postAPi.postService.Delete(id)

	return nil, err
}

func RegisterPostApi(router *mux.Router) *mux.Router {
	post := PostAPI{
		postService: services.NewPostServices(repository.PostRepo),
	}
	router.Methods("GET").Path("/v0/posts/ids").HandlerFunc(post.BuildFuncApi(post.GetIDs))
	router.Methods("POST").Path("/v0/posts").HandlerFunc(post.BuildFuncApi(post.Create))
	router.Methods(http.MethodPut).Path(fmt.Sprintf("/v0/posts/{%s}", POST_ID)).HandlerFunc(post.BuildFuncApi(post.Update))
	router.Methods(http.MethodDelete).Path(fmt.Sprintf("/v0/posts/{%s}", POST_ID)).HandlerFunc(post.BuildFuncApi(post.Delete))

	return router
}
