package main

import (
  //"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"testing"
	"strings"
  "github.com/garyburd/redigo/redis"

  "external"
  "model"
)

//integration test, needs to start the http server and redis server first
func TestIntegrationSetup(t *testing.T) {
  cleanRedisRelated()
}

func TestCreateLocation(t *testing.T) {
  request_body := `{"location":"shanghai"}`
  reader := strings.NewReader(request_body)
  resp, err := http.Post("http://localhost:8080/location", "application/json", reader)
  if resp.StatusCode != 201 {
    t.Fatal("response code should be 201, but actually is ", resp.StatusCode)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if (err != nil) { 
  	t.Fatal(err)
  }
  hash := map[string]string{}
  err = json.Unmarshal(body, &hash)
  if (err != nil) { 
  		t.Fatal(err)
  }
	if hash["location"] != "shanghai" {
		t.Fatal("response is not correct: ", string(body))
	}
	if len(hash) != 1 {
	 t.Fatal("response key number is not 1: ", string(body))
	}
}

func TestCreateLocationDuplicate(t *testing.T) {
  request_body := `{"location":"shanghai"}`
  reader := strings.NewReader(request_body)
  resp, err := http.Post("http://localhost:8080/location", "application/json", reader)
  if resp.StatusCode != 409 {
    t.Fatal("response code should be 409, but actually is ", resp.StatusCode)
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if (err != nil) { 
    t.Fatal(err)
  }
  hash := map[string]string{}
  err = json.Unmarshal(body, &hash)
  if (err != nil) { 
      t.Fatal(err)
    }
  if hash["error"] != "Name already exists" {
    t.Fatal("response is not correct: ", string(body))
  }
  if len(hash) != 1 {
   t.Fatal("response key number is not 1: ", string(body))
  }
}

func TestGetList(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/location")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  	if (err != nil) { 
  		t.Fatal(err)
  	}
  	hash := []string{}
  	err = json.Unmarshal(body, &hash)
	if (err != nil) { 
  		t.Fatal(err)
  	}
  	if len(hash) != 1 {
		t.Fatal("response key number is not 1: ", string(body))
  	}
  	if hash[0] != "shanghai" {
		t.Fatal("response key number is not correct: ", string(body))
  	}
}

func TestGetWeather(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/location/shanghai")
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  	if (err != nil) { 
  		t.Fatal(err)
  	}
  	hash := map[string]interface{} {}
  	err = json.Unmarshal(body, &hash)
	if (err != nil) { 
  		t.Fatal(err)
  	}
  	if len(hash) != 1 {
		t.Fatal("response key number is not 1: ", string(body))
  	}
  	if _, ok := hash["weather"]; !ok {
		t.Fatal("response key doesn't contain weather ", string(body))
  	}
}

func TestDeleteLocation(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8080/location/shanghai", nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  	if (err != nil) { 
  		t.Fatal(err)
  	}
  	if len(body) != 0 {
  		t.Fatal("body is not nil: ", body)
  	}
}

func TestIntegrationTeardown(t *testing.T) {
  cleanRedisRelated()
}

func cleanRedisRelated() {
  conn, _ := redis.DialURL(external.GetConfig().RedisURI)

  //log.Printf("model.GetLocationList():%v", model.GetLocationList())
  for _, name := range model.GetLocationList() {
    conn.Do("DEL", name)
  }
  conn.Do("DEL", model.Set_Key)
}