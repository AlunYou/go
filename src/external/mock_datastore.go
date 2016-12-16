package external

type MockDataStore struct{
	s map[string]string
	m map[string]string
}

func NewMockDataStore(uri string) *MockDataStore {
	store := new(MockDataStore)
	store.s = make(map[string]string)
	store.m = make(map[string]string)
	return store
}

func (ds *MockDataStore) Set_GetKeys(set string) []string {
	keys := make([]string, 0, len(ds.s)) 
	for k := range ds.s {
		keys = append(keys, k)
	}
	return keys
}

func (ds *MockDataStore) Set_AddKey(set string, key string) bool {
	ds.s[key] = ""
	return true
}

func (ds *MockDataStore) Set_DelKey(set string, key string) bool {
	delete(ds.s, key)
	return true
}

func (ds *MockDataStore) Set_IsKey(set string, key string) (bool, error) {
	_, ok := ds.s[key]
	if ok {
		return true, nil
	} else {
		return false, nil
	}
}


func (ds *MockDataStore) Str_Write(key string, value string, expire uint32) bool {
	ds.m[key] = value
	return true
}

func (ds *MockDataStore) Str_Read(key string) (string, error) {
	return ds.m[key], nil
}

func (ds *MockDataStore) Key_Del(key string) (bool, error) {
	delete(ds.m, key)
	return true, nil
}