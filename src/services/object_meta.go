package services

import (
	"dongpham/constant"
	"dongpham/repository"
	"dongpham/utils"
	"errors"
)

const ID_FIELDS = "id"

type ObjectMetaServices struct {
}

var metaObjectGetter = map[string]func(oids []string, fields []string) (map[string]interface{}, error){}

func (service *ObjectMetaServices) GetObjectMetaByIDsAndFields(oids []string, fields []string) (map[string]interface{}, error) {
	mapObjectWithType := map[string][]string{}

	for _, oid := range oids {
		oType := utils.GetObjectTypeFromOID(oid)
		if _, ok := mapObjectWithType[oType]; ok {
			mapObjectWithType[oType] = []string{}
		}
		mapObjectWithType[oType] = append(mapObjectWithType[oType], oid)
	}
	result := map[string]interface{}{}
	for oType, oids := range mapObjectWithType {
		if metaObjectGetter[oType] == nil {
			return nil, errors.New("oType is invalid %" + oType)
		}
		metaObjects, err := metaObjectGetter[oType](oids, fields)
		if err != nil {
			return nil, err
		}

		for key, val := range metaObjects {
			result[key] = val
		}

	}

	return result, nil
}

func NewObjectMetaServices() *ObjectMetaServices {
	return &ObjectMetaServices{}
}

func init() {
	postService := NewPostServices(repository.PostRepo)
	metaObjectGetter[constant.META_OBJECT_POST] = func(oids []string, fields []string) (i map[string]interface{}, e error) {
		if len(oids) == 0 || len(fields) == 0 {
			return nil, nil
		}
		oType := utils.GetObjectTypeFromOID(oids[0])
		posts, err := postService.GetPostByIDs(utils.ParseOIDsToID(oids), fields)
		if err != nil {
			return nil, err
		}
		result := map[string]interface{}{}
		for key, post := range posts {
			if post != nil {
				result[utils.GenerateOID(oType, key)] = post
			}
		}
		return result, err
	}
}

