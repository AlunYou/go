package external

import (
	"net/http"
	"github.com/garyburd/redigo/redis"
)


type Config struct {
	WeatherWervice WeatherService
	WeatherAppId string
	WeatherApiEndpoint string
	RedisURI string
	CacheTimeout uint32
}

func NewConfig() *Config {
	config := new(Config)
	config.RedisURI = "redis://localhost:6379"
	config.WeatherWervice = &OpenWeatherService{}
	config.WeatherAppId = ""
	config.WeatherApiEndpoint = "http://api.openweathermap.org/data/2.5/weather?q=%s&APPID=%s" 
	config.CacheTimeout = 60 * 60 // 1 hour
	return config
}

//implement config singleton. The config initialization happens when just enter main(), so no need to be thread safe
var (
	config *Config
	redis_go *RedisGo
	redis_pool *redis.Pool
)

func GetConfig() *Config {
	if config == nil {
		config = NewConfig()
	}
	return config
}

func GetRedisGo() *RedisGo {
	if redis_go == nil {
		redis_go = &RedisGo{DialURL: redis.DialURL, Strings: redis.Strings, String: redis.String, Bool: redis.Bool}
	}
	return redis_go
}

func DestroyConfig() {
	if config != nil {
		config = nil
	}
}

var GetDataStore = func() DataStore {
	if redis_pool == nil {
		redis_pool = &redis.Pool{
	        MaxIdle:   80,
	        MaxActive: 1200, // max number of connections
	        Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(GetConfig().RedisURI)
			    if err != nil {
			        panic(err.Error())
			    }
			    return c, err
			} }	
	}
	conn := redis_pool.Get()
	return NewRedisDataStore(conn)
}

func UseRedisDataStoreWithMockConn() {
	var store DataStore
	GetDataStore = func() DataStore {
		if store == nil {
			conn, _ := GetRedisGo().DialURL("")
		    store = NewRedisDataStore(conn)
		}
		return store
	}
}

func UseMockDataStore() {
	var store DataStore
	GetDataStore = func() DataStore {
		if store == nil {
			store = NewMockDataStore("")
		}
		return store
	}
}

// http client is concurrency safe, so use singleton
var GetHttpClient = func() *http.Client {
	var client *http.Client
	if client == nil {
		client = &http.Client{}
	}
	return client
}