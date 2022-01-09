package redis

const CACHE_OBJECT_DATA = "cacheObjectData"

func CacheWithNIDs(ids []uint32,fields []string, getter func(ids []uint32) map[string]interface{}) map[uint32]string {
	return nil
}

func CacheWithIDs(oids []string, getter func(ids []uint32) map[string]interface{}) map[uint32]string {
	return nil
}