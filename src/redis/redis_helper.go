package redis

func CacheWithNIDs(ids []uint32, getter func(ids []uint32) map[string]interface{}) map[uint32]string {
	return nil
}

func CacheWithIDs(ids []string, getter func(ids []uint32) map[string]interface{}) map[uint32]string {
	return nil
}