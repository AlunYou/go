package model

import (
	"external"
)

const Set_Key = "location_set_key"

func GetLocationList() []string {
	ds := external.GetDataStore()
	list := ds.Set_GetKeys(Set_Key)
	return list
}

func CreateLocation(name string) bool {
	ds := external.GetDataStore()
	return ds.Set_AddKey(Set_Key, name)
}

func FindLocation(name string) (bool, error) {
	ds := external.GetDataStore()
	is, err := ds.Set_IsKey(Set_Key, name)
	if err != nil {
		return false, err
	}
	return is, nil
}

func DeleteLocation(name string) bool {
	ds := external.GetDataStore()
	res := ds.Set_DelKey(Set_Key, name)
	ds.Key_Del(name)
	return res
}

func GetWeather(name string) (string, error) {
	ds := external.GetDataStore()
	//first find the cache
	v, err := ds.Str_Read(name)
	if err != nil {
		return "", err
	}
	if v == "" {
		weather,err := external.GetConfig().WeatherWervice.Lookup(name)
		if err != nil {
			return "", err
		}
		ds.Str_Write(name, weather, external.GetConfig().CacheTimeout)
		return weather, nil
	} else {
		return v, nil
	}
}

