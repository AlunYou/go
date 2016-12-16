package external

import (
	"testing"
	"net/http"

	"test_util"
)

func setup() {
	DestroyConfig()
	GetConfig()
}
func teardown() {
}

func Test_WeatherService(t *testing.T) {
	setup()
	defer teardown()

	//mock a http server with correct response
	server, client := test_util.MockHttpServer(200, `{"coord":{"lon":121.46,"lat":31.22},"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01d"}],"base":"stations","main":{"temp":288.57,"pressure":1014,"humidity":45,"temp_min":287.15,"temp_max":290.15},"visibility":10000,"wind":{"speed":5,"deg":110},"clouds":{"all":0},"dt":1459762200,"sys":{"type":1,"id":7452,"message":0.0206,"country":"CN","sunrise":1459719525,"sunset":1459764940},"id":1796236,"name":"Shanghai","cod":200}`)
	defer server.Close()

	//set the client of app, so that it will connect to the mock server
	GetHttpClient = func() *http.Client{
		return client
	}

	//start testing
    json, _ := GetConfig().WeatherWervice.Lookup("shanghai")
    if json != `{"weather":[{"description":"clear sky","icon":"01d","id":800,"main":"Clear"}]}` {
    	t.Fatal("generate wrong json: ", json)
    }
}
