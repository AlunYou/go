package external

import (
	"testing"
	"github.com/garyburd/redigo/redis"

	"test_util"
)

func setupRedisTest(t *testing.T) {
	DestroyConfig()

	UseRedisDataStoreWithMockConn()

	GetRedisGo().DialURL = func (rawurl string, options ...redis.DialOption) (redis.Conn, error) { 
	  	command_list := []string{}
	  	
	  	return &MockRedisConn{command_list: command_list}, nil
  	}
  	GetRedisGo().Strings = func(reply interface{}, err error) ([]string, error) { 
	  	return nil, nil
  	}
  	GetRedisGo().String = func(reply interface{}, err error) (string, error) { 
	  	return "", nil
  	}
  	GetRedisGo().Bool = func(reply interface{}, err error) (bool, error) { 
	  	return false, nil
  	}
  	GetRedisGo().Close = func() bool { 
	  	return false
  	}
}

func teardownRedisTest() {
	DestroyConfig()
}

func TestSetAddKey(t *testing.T) {
	setupRedisTest(t)
	defer teardownRedisTest()

	store := GetDataStore()
	store.Set_AddKey("test_key", "shanghai")
	command_list := store.(*RedisDataStore).Conn.(*MockRedisConn).command_list
	if !test_util.CompareSlice(command_list, []string{"SADD test_key shanghai"}) {
		t.Fatal("command not match, actual:", command_list)
	}
}

func TestSetDelKey(t *testing.T) {
	setupRedisTest(t)
	defer teardownRedisTest()

	store := GetDataStore()
	store.Set_DelKey("test_key", "shanghai")
	command_list := store.(*RedisDataStore).Conn.(*MockRedisConn).command_list
	if !test_util.CompareSlice(command_list, []string{"SREM test_key shanghai"}) {
		t.Fatal("command not match, actual:", command_list)
	}
}

func TestSetGetKeys(t *testing.T) {
	setupRedisTest(t)
	defer teardownRedisTest()

	store := GetDataStore()
	store.Set_GetKeys("test_key")
	command_list := store.(*RedisDataStore).Conn.(*MockRedisConn).command_list
	if !test_util.CompareSlice(command_list, []string{"SMEMBERS test_key"}) {
		t.Fatal("command not match, actual:", command_list)
	}
}

func TestSetIsKey(t *testing.T) {
	setupRedisTest(t)
	defer teardownRedisTest()

	store := GetDataStore()
	store.Set_IsKey("test_key", "shanghai")
	command_list := store.(*RedisDataStore).Conn.(*MockRedisConn).command_list
	if !test_util.CompareSlice(command_list, []string{"SISMEMBER test_key shanghai"}) {
		t.Fatal("command not match, actual:", command_list)
	}
}

func TestStringWrite(t *testing.T) {
	setupRedisTest(t)
	defer teardownRedisTest()

	store := GetDataStore()
	store.Str_Write("test_key", "shanghai", 100)
	command_list := store.(*RedisDataStore).Conn.(*MockRedisConn).command_list
	if !test_util.CompareSlice(command_list, []string{"MULTI", "SET test_key shanghai", "EXPIRE test_key 100", "EXEC"}) {
		t.Fatal("command not match, actual:", command_list)
	}
}

func TestStringRead(t *testing.T) {
	setupRedisTest(t)
	defer teardownRedisTest()

	store := GetDataStore()
	store.Str_Read("test_key")
	command_list := store.(*RedisDataStore).Conn.(*MockRedisConn).command_list
	if !test_util.CompareSlice(command_list, []string{"GET test_key"}) {
		t.Fatal("command not match, actual:", command_list)
	}
}

func TestKeyDel(t *testing.T) {
	setupRedisTest(t)
	defer teardownRedisTest()

	store := GetDataStore()
	store.Key_Del("test_key")
	command_list := store.(*RedisDataStore).Conn.(*MockRedisConn).command_list
	if !test_util.CompareSlice(command_list, []string{"DEL test_key"}) {
		t.Fatal("command not match, actual:", command_list)
	}
}

