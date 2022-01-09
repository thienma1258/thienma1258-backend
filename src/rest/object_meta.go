package rest

import (
	"dongpham/model"
	"dongpham/services"
	"github.com/gorilla/mux"
	"net/http"
)

type ObjectMetaAPI struct {
	metaService *services.ObjectMetaServices
	Api
}

type ObjectMetaRequest struct {
	OIDs   []string `json:"oids"`
	Fields []string `json:"fields"`
}

func (metaAPI *ObjectMetaAPI) GetIDs(request *model.ApiRequest) (interface{}, error) {
	queryRequest := &ObjectMetaRequest{}
	err := json.NewDecoder(request.Body).Decode(&queryRequest)
	if err != nil {
		return nil, err
	}
	ids, err := metaAPI.metaService.GetObjectMetaByIDsAndFields(queryRequest.OIDs, queryRequest.Fields)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func RegisterObjectMetaApi(router *mux.Router) *mux.Router {
	post := ObjectMetaAPI{
		metaService: services.NewObjectMetaServices(),
	}
	router.Methods(http.MethodPut).Path("/v0/objects").HandlerFunc(post.BuildFuncApi(post.GetIDs))

	return router
}
