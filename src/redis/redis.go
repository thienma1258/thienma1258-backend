package redis

import (
	"context"
	localConfig "dongpham/config"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

const NUM_PIPELINE_PROCESS = 100

var client *redis.Client

var luaSHA string
var luaLock sync.RWMutex
var config *RedisConnectionSetting

const luaScript = `
	local results = {}
	for i = 1, table.getn(KEYS) do
		results[i] = redis.call("HMGET", KEYS[i], unpack(ARGV))
	end
	return results
`

type RedisConnectionSetting struct {
	Ctx     context.Context
	Address string
	DB      int
	Timeout time.Duration
}

func MSet(cKey string, fields map[string]interface{}) {
	if len(fields) == 0 {
		debug.PrintStack()
		log.Printf("HMSet with empty fields")
	}
	client := getClient()
	var pairs []interface{}

	for key, value := range fields {
		pairs = append(pairs, key, value)

	}
	result := client.MSet(config.Ctx, fields)
	if result.Err() != nil {
		log.Printf("HMSet error %v", result.Err())
	}
}

func Delete(cKeys ...string) {
	client := getClient()
	client.Del(config.Ctx, cKeys...)

}

func RegisterRedisConnection(setting *RedisConnectionSetting) {
	client = redis.NewClient(&redis.Options{
		Addr:     setting.Address,
		Password: "",         // no password set
		DB:       setting.DB, // use default DB
	})
}

func getClient() *redis.Client {
	return client
}

func MGet(conn int, cKeys []string) map[string]string {
	client := getClient()
	result := client.MGet(config.Ctx, cKeys...)
	if result.Err() != nil {
		return nil
	}
	results := make(map[string]string)
	values := result.Val()
	for i := 0; i < len(cKeys); i++ {
		ckey := cKeys[i]
		if values[i] == nil {
			results[ckey] = ""
		} else {
			results[ckey] = fmt.Sprintf("%s", values[i])
		}
	}
	return results
}

func HGet(cKey string, field string) (string, error) {
	client := getClient()
	result := client.HGet(config.Ctx, cKey, field)
	if result.Err() != nil && !strings.Contains(result.Err().Error(),"nil") {
		log.Printf("HGet error %v", result.Err())
		return "", result.Err()
	}
	return result.Val(), nil
}

func HMGet( cKey string, cFields ...string) map[string]string {
	client := getClient()
	result := client.HMGet(config.Ctx,cKey, cFields...)
	if result.Err() != nil {
		return nil
	}
	results := make(map[string]string)
	values := result.Val()
	// result.String()
	for i := 0; i < len(cFields); i++ {
		cField := cFields[i]
		if values[i] == nil {
			results[cField] = ""
		} else {
			results[cField] = fmt.Sprintf("%s", values[i])
		}
	}
	return results
}

func HSetNX(cKey string, field string, val []byte) ( error) {
	client := getClient()
	result := client.HSetNX(config.Ctx, cKey, field, val)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func SetTimeout(ctx context.Context, cKey string, value interface{}, timeout time.Duration) {
	client := getClient()
	client.Set(ctx, cKey, value, timeout)
}

func Set(ctx context.Context, conn int, cKey string, value interface{}) {
	// log.Printf("SetZ conn=%d key=%s", conn, cKey)
	client := getClient()
	client.Set(ctx, cKey, value, config.Timeout)
}

func HMSet(ctx context.Context, cKey string, fields map[string]interface{}) {
	if len(fields) == 0 {
		debug.PrintStack()
		log.Printf("HMSet with empty fields")
	}
	client := getClient()
	result := client.HMSet(ctx, cKey, fields)
	if result.Err() != nil {
		log.Printf("HMSet error %v", result.Err())
	}

	applyTimeout(cKey)
}

// AcquireLock acquire lock for syncing multiple processes:
// + Return empty lock value if it cannot acquire the lock.
// + Return a lock value for release lock later.
func AcquireLock(
	conn int,
	lockKey string,
	wait bool,
	lockTime time.Duration,
	waitTime time.Duration,
	waitTimeOut time.Duration,
) string {
	client := getClient()

	cKey := "rlock:" + lockKey
	nowTs := time.Now().Nanosecond()
	lockValue := fmt.Sprintf("%d-%d", nowTs, rand.Int()) // #nosec
	success := client.SetNX(config.Ctx, cKey, lockValue, lockTime).Val()
	if success {
		return lockValue
	}
	if !wait {
		return ""
	}

	timeOutTs := nowTs + int(waitTimeOut.Nanoseconds())

	for {
		success = client.SetNX(config.Ctx, cKey, lockValue, lockTime).Val()
		if success {
			return lockValue
		}
		if time.Now().Nanosecond() > timeOutTs {
			return ""
		}
		time.Sleep(waitTime/2 + time.Duration(rand.Intn(int(waitTime/2)))) // #nosec
	}
}

// ReleaseLock release the lock. lockValue is the value received from AcquireLock method.
func ReleaseLock(conn int, lockKey string, lockValue string) {
	client := getClient()

	cKey := "rlock:" + lockKey
	cVal := client.Get(config.Ctx, cKey).Val()
	if cVal == lockValue {
		client.Del(config.Ctx, cKey)
	}
}

func applyTimeout(cKey string) {
	timeout := config.Timeout
	if timeout > 0 {
		client.Expire(config.Ctx, cKey, timeout)
	}
}

func HGetMultipleFieldsLuaScript(
	cKeys []string, fields []string,
) *map[string](map[string]([]byte)) {
	client := getClient()
	var err error
	luaLock.RLock()
	sha := luaSHA
	luaLock.RUnlock()

	if sha == "" {
		sha, err = client.ScriptLoad(config.Ctx, luaScript).Result()
		if err != nil {
			log.Printf("Error while loading script %v\n", err)
			return nil
		}
		luaLock.Lock()
		luaSHA = sha
		luaLock.Unlock()
	}

	argv := make([]interface{}, len(fields))
	for i := 0; i < len(fields); i++ {
		argv[i] = fields[i]
	}

	results, err := client.EvalSha(config.Ctx, sha, cKeys, argv...).Result()
	if err != nil {
		log.Printf("Error while load from cache %v\n", err)
		return nil
	}

	values := results.([]interface{})
	total := len(values)
	totalFields := len(fields)
	items := make(map[string](map[string]([]byte)))

	for i := 0; i < total; i++ {
		item := make(map[string]([]byte))

		val := values[i].([]interface{})
		for j := 0; j < totalFields; j++ {
			field := fields[j]
			switch v := val[j].(type) {
			case int:
				item[field] = val[j].([]byte)
			case string:
				item[field] = []byte(val[j].(string))
			case nil:
				item[field] = nil
			default:
				fmt.Printf("I don't know about type %T!\n", v)
			}
		}
		items[cKeys[i]] = item
	}
	return &items
}

func HMSetMultipleKeys(items map[string](map[string]interface{})) {
	client := getClient()
	timeout := config.Timeout
	pipe := client.Pipeline()
	for cKey, values := range items {
		pipe.HMSet(config.Ctx, cKey, values)
		if timeout > 0 {
			pipe.Expire(config.Ctx, cKey, timeout)
		}
	}
	_, err := pipe.Exec(config.Ctx)
	if err != nil {
		log.Printf("HMSetMultipleKeys error=%v", err)
	}
}

func init() {
	config = &RedisConnectionSetting{
		Ctx:     context.Background(),
		Address: localConfig.RedisAddr,
		DB:      0,
		Timeout: 0,
	}
	RegisterRedisConnection(config)

}
