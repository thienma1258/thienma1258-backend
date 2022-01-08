package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"math/rand"
	"runtime/debug"
	"sync"
	"time"
)

const CACHE_OBJECT_DATA = "cacheObjectData"
const NUM_PIPELINE_PROCESS = 100

var client *redis.Client

var luaSHA = make(map[int]string)
var luaLock sync.RWMutex
var config *RedisConnectionSetting

type RedisConnectionSetting struct {
	Ctx     context.Context
	Address string
	DB      int
	Timeout time.Duration
	client  *redis.Client
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

func Delete(ctx context.Context, cKeys ...string) {
	client := getClient()
	client.Del(ctx, cKeys...)

}

func RegisterRedisConnection(connection int, setting *RedisConnectionSetting) {
	conn := redis.NewClient(&redis.Options{
		Addr:     setting.Address,
		Password: "",         // no password set
		DB:       setting.DB, // use default DB
	})
	setting.client = conn
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
