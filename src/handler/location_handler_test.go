package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"strings"
	"encoding/json"
	"html"
	
	"model"
	"external"
)

const ATTACK_STR = `<script> alert('')</script>`

func setup() {
  external.UseMockDataStore()
}

func teardown() {
}

func TestHandleCreateLocationOk(t *testing.T) {
	setup()
	defer teardown()

	testCreateOk(t, "shanghai", "shanghai")
	testCreateOk(t, ATTACK_STR, html.EscapeString(ATTACK_STR))
}

func testCreateOk(t *testing.T, location_input string, location_output string) {
	request_body := fmt.Sprintf(`{"location":"%s"}`, location_input)
    reader := strings.NewReader(request_body)
	req, err := http.NewRequest("POST", "http://localhost:8080/location", reader)
	if err != nil {
	    t.Fatal(err)
	}

	w := httptest.NewRecorder()
	HandleCreateLocation(w, req)
	if w.Code != 201 {
		t.Fatal("response code should be 201, but actually is ", w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
	if (err != nil) { 
		t.Fatal(err)
	}
	hash := map[string]string{}
	err = json.Unmarshal(body, &hash)
	if (err != nil) { 
		t.Fatal(err, string(body))
	}
	if hash["location"] != location_output {
		t.Fatal("response is not correct: ", string(body))
	}
	if len(hash) != 1 {
	 t.Fatal("response key number is not 1: ", string(body))
	}
}

//also test case insensitive
func TestHandleCreateLocationConflict(t *testing.T) {
	setup()
	model.CreateLocation("shanghai")
	defer teardown()

	testCreateConflict(t, "Shanghai")
	testCreateConflict(t, "shanghai")
}

func testCreateConflict(t *testing.T, location string) {
	request_body := fmt.Sprintf("{\"location\":\"%s\"}", location)
    reader := strings.NewReader(request_body)
	req, err := http.NewRequest("POST", "http://localhost:8080/location", reader)
	if err != nil {
	    t.Fatal(err)
	}

	w := httptest.NewRecorder()
	HandleCreateLocation(w, req)
	if w.Code != 409 {
		t.Fatal("response code should be 409, but actually is ", w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
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

func TestHandleGetListOk(t *testing.T) {
	setup()
	model.CreateLocation("shanghai")
	defer teardown()

	req, err := http.NewRequest("GET", "http://localhost:8080/location", nil)
	if err != nil {
	    t.Fatal(err)
	}

	w := httptest.NewRecorder()
	HandleLocationList(w, req)
	if w.Code != 200 {
		t.Fatal("response code should be 200, but actually is ", w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
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

//also test case insensitive
func TestHandleDeleteLocationOK(t *testing.T) {
	setup()
	model.CreateLocation("shanghai")
	model.CreateLocation("beijing")
	model.CreateLocation(html.EscapeString(ATTACK_STR))
	defer teardown()

	testDeleteLocation(t, "shanghai")
	testDeleteLocation(t, "BeiJing")
	testDeleteLocation(t, ATTACK_STR)
}

func testDeleteLocation(t *testing.T, location string) {
	url := fmt.Sprintf("http://localhost:8080/location/%s", location)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
	    t.Fatal(err)
	}

	w := httptest.NewRecorder()
	HandleDeleteLocation(w, req)
	if w.Code != 200 {
		t.Fatal("response code should be 200, but actually is ", w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
	if (err != nil) { 
		t.Fatal(err)
	}

	if len(body) != 0 {
		t.Fatal("response body should be nil: ", string(body))
  	}
}

func TestHandleDeleteLocationNotFound(t *testing.T) {
	setup()
	model.CreateLocation("shanghai")
	defer teardown()

	req, err := http.NewRequest("DELETE", "http://localhost:8080/location/beijing", nil)
	if err != nil {
	    t.Fatal(err)
	}

	w := httptest.NewRecorder()
	HandleDeleteLocation(w, req)
	if w.Code != 404 {
		t.Fatal("response code should be 404, but actually is ", w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
	if (err != nil) { 
		t.Fatal(err)
	}
	hash := map[string]string{}
	err = json.Unmarshal(body, &hash)
	if (err != nil) { 
		t.Fatal(err)
	}
	if hash["error"] != "Name not found" {
		t.Fatal("response is not correct: ", string(body))
	}
	if len(hash) != 1 {
	 	t.Fatal("response key number is not 1: ", string(body))
	}
}

//also test case insensitive
func TestHandleGetLocationOK(t *testing.T) {
	setup()
	model.CreateLocation("shanghai")
	model.CreateLocation(html.EscapeString(ATTACK_STR))
	defer teardown()

	testGetLocation(t, "Shanghai")
	testGetLocation(t, "shanghai")
	testGetLocation(t, ATTACK_STR)
}

func testGetLocation(t *testing.T, location string) {
	url := fmt.Sprintf("http://localhost:8080/location/%s", location)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
	    t.Fatal(err)
	}

	w := httptest.NewRecorder()
	HandleGetLocation(w, req)
	if w.Code != 200 {
		t.Fatal("response code should be 200, but actually is ", w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
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

func TestHandleGetLocationNotFound(t *testing.T) {
	setup()
	model.CreateLocation("shanghai")
	defer teardown()

	req, err := http.NewRequest("DELETE", "http://localhost:8080/location/beijing", nil)
	if err != nil {
	    t.Fatal(err)
	}

	w := httptest.NewRecorder()
	HandleGetLocation(w, req)
	if w.Code != 404 {
		t.Fatal("response code should be 404, but actually is ", w.Code)
	}

	body, err := ioutil.ReadAll(w.Body)
	if (err != nil) { 
		t.Fatal(err)
	}
	hash := map[string]string{}
	err = json.Unmarshal(body, &hash)
	if (err != nil) { 
		t.Fatal(err)
	}
	if hash["error"] != "Name not found" {
		t.Fatal("response is not correct: ", string(body))
	}
	if len(hash) != 1 {
	 	t.Fatal("response key number is not 1: ", string(body))
	}
}

