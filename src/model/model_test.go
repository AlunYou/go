package model

import (
	"testing"

	"external"
	"test_util"
)
const GOOD_WEATHER = `{"weather" : "bad"}`
const BAD_WEATHER = `{"weather" : "bad"}`

type MockWeatherService struct {
}

func (service *MockWeatherService) Lookup(location string) (string, error) {
	return BAD_WEATHER, nil
}

func setup() {
	external.UseMockDataStore()
	external.GetConfig().WeatherWervice = &MockWeatherService{}
}

func teardown() {

}

func TestLocation(t *testing.T) {
	setup()
	defer teardown()

	//test create 
	res := CreateLocation("shanghai")
	if !res {
		t.Fatal("create location return false")
	}
	keys := external.GetDataStore().Set_GetKeys(Set_Key)
	if !test_util.CompareSlice(keys, []string{"shanghai"}) {
		t.Fatal("keys not match, actual:", keys)
	}

	//test get list
	location_list := GetLocationList()
	if !test_util.CompareSlice(location_list, []string{"shanghai"}) {
		t.Fatal("keys not match, actual:", location_list)
	}

    //test delete
	res = DeleteLocation("shanghai")
	if !res {
		t.Fatal("delete location return false")
	}
	location_list = GetLocationList()
	if !test_util.CompareSlice(location_list, []string{}) {
		t.Fatal("keys not match after delete, actual:", location_list)
	}
}

func TestGetWeatherFromCache(t *testing.T) {
	setup()
	CreateLocation("shanghai")
	external.GetDataStore().Str_Write("shanghai", GOOD_WEATHER, 100)
	defer teardown()

	weather, _ := GetWeather("shanghai")
	if weather != GOOD_WEATHER {
		t.Fatal("should return cached weather")
	}
}

func TestGetWeatherFromWeatherService(t *testing.T) {
	setup()
	CreateLocation("shanghai")
	defer teardown()

	weather, _ := GetWeather("shanghai")
	if weather != BAD_WEATHER {
		t.Fatal("should return weather from external service")
	}

	cached, _ := external.GetDataStore().Str_Read("shanghai")
	if cached != BAD_WEATHER {
		t.Fatal("should cache weather to datastore")
	}
}