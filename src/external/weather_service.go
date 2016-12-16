package external

import (
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
)

const DUMMYWEATHER = `{"weather":""}`

type WeatherService interface {
	Lookup(location string) (string, error)
}

type OpenWeatherService struct {

}

func (service *OpenWeatherService) Lookup(location string) (string, error) {
	app := GetConfig()
	url := fmt.Sprintf(app.WeatherApiEndpoint, location, app.WeatherAppId)
	client := GetHttpClient()
	resp, err := client.Get(url)
	if(err != nil){
		log.Printf("lookup error: %s", err)
		return "", err
	}
	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  	if (err != nil) { 
  		return "", err 
  	}
	log.Printf("url: %s, code: , response: %s", url, resp.StatusCode, string(body));
	if resp.StatusCode != 200 {
		return DUMMYWEATHER, nil
	}
  	return gen_json(body)
}

func gen_json(body []byte) (string, error) {
  	var hash interface{}
	err := json.Unmarshal(body, &hash)
	if (err != nil) { 
  		return "", err 
  	}

	asserted_hash, ok := hash.(map[string]interface{})
	if !ok {
		return DUMMYWEATHER, fmt.Errorf("api not return a hash")
	}
	assert_weather, ok := asserted_hash["weather"]
	if !ok {
		return DUMMYWEATHER, nil
		//return "", fmt.Errorf("api not contain weather key")
	}
	bytes, err := json.Marshal(assert_weather)
	if (err != nil) { 
  		return DUMMYWEATHER, err 
  	}

	weather_body := map[string]interface{} {}
	weather_body["weather"] = assert_weather
	bytes, err = json.Marshal(weather_body)
	if (err != nil) { 
  		return DUMMYWEATHER, err 
  	}
	log.Printf("weather_body: %v", string(bytes));
	return string(bytes), err
}