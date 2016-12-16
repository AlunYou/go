package external

import (
	"log"
	"github.com/garyburd/redigo/redis"
)

type RedisGo struct {
	DialURL func(rawurl string, options ...redis.DialOption) (redis.Conn, error)
	Strings func(reply interface{}, err error) ([]string, error)
	String func(reply interface{}, err error) (string, error) 
	Bool func(reply interface{}, err error) (bool, error)
	Close func() bool
}

type DataStore interface {
	//set operation
	Set_GetKeys(set string) []string
    Set_AddKey(set string, key string) bool
    Set_DelKey(set string, key string) bool
    Set_IsKey(set string, key string) (bool, error)

    //string operation
    Str_Write(key string, value string, expire uint32) bool
    Str_Read(key string) (string, error)

    //key operation
    Key_Del(key string) (bool, error)
}

type RedisDataStore struct{
	Conn redis.Conn
}

func NewRedisDataStore(conn redis.Conn) *RedisDataStore {
	store := new(RedisDataStore)
	store.Conn= conn //GetRedisGo().DialURL(uri)
	return store
}

func (ds *RedisDataStore) Set_GetKeys(set string) []string {
	resp, err := ds.Conn.Do("SMEMBERS", set)
	if err != nil {
		log.Printf("Set_GetKeys error: %s", err)
	    return nil
	}
	if v, ok := resp.([]interface{}); ok {
		ary, err := GetRedisGo().Strings(v, err)
		if err != nil {
			log.Printf("Set_GetKeys error: %s", err)
		    return nil
		}
		return ary
	} else {
		log.Printf("Set_GetKeys error: %s", "return wrong type")
		return nil
	}
}

func (ds *RedisDataStore) Set_AddKey(set string, key string) bool {
	_, err := ds.Conn.Do("SADD", set, key)
	if err != nil {
		log.Printf("Set_AddKey error: %s", err)
	    return false
	}
	return true
}

func (ds *RedisDataStore) Set_DelKey(set string, key string) bool {
	_, err := ds.Conn.Do("SREM", set, key)
	if err != nil {
		log.Printf("Set_DelKey error: %s", err)
	    return false
	}
	return true
}

func (ds *RedisDataStore) Set_IsKey(set string, key string) (bool, error) {
	resp, err := ds.Conn.Do("SISMEMBER", set, key)
	if err != nil {
		log.Printf("Set_IsKey error: %s", err)
	    return false, err
	}
	v, err := GetRedisGo().Bool(resp, err)
	if err != nil {
		log.Printf("Set_IsKey error: %s", err)
	    return false, err
	} else {
		return v, nil
	}
}

func (ds *RedisDataStore) Str_Write(key string, value string, expire uint32) bool {
	ds.Conn.Send("MULTI")
	ds.Conn.Send("SET", key, value)
	ds.Conn.Send("EXPIRE", key, expire)
	_, err := ds.Conn.Do("EXEC")
	if err != nil {
		log.Printf("Str_Write error: %s", err)
	    return false
	}
	return true
}

func (ds *RedisDataStore) Str_Read(key string) (string, error) {
	resp, err := ds.Conn.Do("GET", key)
	if err != nil {
		log.Printf("Str_Read error: %s", err)
	    return "", err
	}
	if resp == nil {
		return "", nil
	}
	v, err := GetRedisGo().String(resp, err)
	if err != nil {
		log.Printf("Str_Read error: %s", err)
	    return "", err
	} else {
		return v, nil
	}
}

func (ds *RedisDataStore) Key_Del(key string) (bool, error) {
	resp, err := ds.Conn.Do("DEL", key)
	if err != nil {
		log.Printf("Key_Del error: %s", err)
	    return false, err
	}
	v, err := GetRedisGo().Bool(resp, err)
	if err != nil {
		log.Printf("Key_Del error: %s", err)
	    return false, err
	} else {
		return v, nil
	}
}

