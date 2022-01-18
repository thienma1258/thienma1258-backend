package redis

import (
	"dongpham/utils"
	"github.com/fatih/structs"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const CACHE_OBJECT_DATA = "cacheObjectData"

func CacheWithNIDs(key string, ids []int, fields []string, getter func(ids []int) map[string]interface{}) map[int]string {
	return nil
}

func CacheWithIDs(key string, oids []string, getter func(ids []int) map[string]interface{}) map[string]string {
	return nil
}

func CacheWithKey(key string, queryString string, getter func() (interface{}, error)) (interface{}, error) {
	cache, err := HGet(key, queryString)
	if err != nil {
		return nil, err
	}
	if len(cache) > 0 {
		var result interface{}
		err = json.Unmarshal([]byte(cache), &result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	result, err := getter()
	if err != nil {
		return nil, err
	}

	setCache, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}
	err = HSetNX(key, queryString, setCache)
	if err != nil {
		return nil, err
	}
	return result.(interface{}), nil
}

func GetDataFromCacheWithEntityType(ids []int, entityType string, fields []string,
	getter func(ids []int) (map[int]interface{}, error)) (map[int]interface{}, error) {
	cKeys := utils.ParseIDToOID(ids, entityType)
	numField := len(fields)
	cache := HGetMultipleFieldsLuaScript(cKeys, fields)
	var missingOIDS []string
	result := map[int]interface{}{}
	if len(*cache) > 0 {

		for oid, data := range *cache {
			item := map[string]interface{}{}
			if data != nil {
				for field, val := range data {
					if val == nil {
						break
					} else {
						var valObject interface{}
						err := json.Unmarshal(val, &valObject)
						if err != nil {
							return nil, err
						}
						item[field] = valObject
					}
				}
			}
			if len(item) == numField {
				result[utils.GetIDFromOID(oid)] = item
			} else {
				missingOIDS = append(missingOIDS, oid)

			}
		}
	}
	if len(missingOIDS) == 0 {
		return result, nil
	}
	data, err := getter(utils.ParseOIDsToID(missingOIDS))
	if err != nil {
		return nil, err
	}
	cacheData := map[string]map[string]interface{}{}
	for id, item := range data {
		mapStringValueCache := map[string]interface{}{}
		mapStringVal := map[string]interface{}{}

		itemFields := structs.New(item).Fields()

		for _, field := range itemFields {
			fieldKey := field.Tag("json")
			byteVal, err := json.Marshal(field.Value())
			if err != nil {
				return nil, err
			}
			mapStringValueCache[fieldKey] = byteVal
			mapStringVal[fieldKey] = field.Value()
		}

		cacheData[utils.GenerateOID(entityType, id)] = mapStringValueCache
		result[id] = mapStringVal
	}
	HMSetMultipleKeys(cacheData)

	return result, nil

}

func DeleteWrapperCache(key string, deleteFunc func() error) error {
	err := deleteFunc()
	if err != nil {
		return err
	}

	Delete(key)
	return nil
}

func CreateWrapperCache(key string, createFunc func() (int, error)) (int, error) {
	result, err := createFunc()
	if err != nil {
		return 0, err
	}

	Delete(key)
	return result, nil
}

func UpdateWrapperCacheWithIDs(cKey string, entityType string, ids []int, _ []string, updateFunc func(ids []int) error) error {

	err := updateFunc(ids)
	if err != nil {
		return err
	}
	cKeys := utils.ParseIDToOID(ids, entityType)
	cKeys = append(cKeys, cKey)
	Delete(cKeys...)
	return nil

}
