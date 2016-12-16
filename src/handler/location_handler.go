package handler

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"html"

	"model"
)

func HandleLocationList(w http.ResponseWriter, r *http.Request) {
	log.Println("list location")
	list := model.GetLocationList()
	renderOK(w, http.StatusOK, list)
}

func HandleCreateLocation(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024))
    if err != nil {
        renderInternalServerError(w, err)
        return
    }
    if err := r.Body.Close(); err != nil {
        renderInternalServerError(w, err)
        return
    }
    location_hash := map[string]string{}
    if err := json.Unmarshal(body, &location_hash); err != nil {
		renderInternalServerError(w, err)
        return
	}

	location := strings.ToLower(html.EscapeString(location_hash["location"]))
	location_hash["location"] = location
	log.Println("add:", location)
	find, err := model.FindLocation(location);
	if err != nil {
		renderInternalServerError(w, err)
        return
	}
	if find {
		hash := map[string]string{ "error": "Name already exists"}
		renderClientError(w, http.StatusConflict, hash)
		return
	}

	log.Printf("request location: %v", location)
	result := model.CreateLocation(location)
	if !result {
		renderInternalServerError(w, fmt.Errorf("can not save location"))
        return
	}
	
	renderOK(w, http.StatusCreated, location_hash)
}

func HandleDeleteLocation(w http.ResponseWriter, r *http.Request) {
	location := html.EscapeString(r.URL.Path[len("/location/"):])
	location = strings.ToLower(location)
	log.Println("delete:", location)
	find, err := model.FindLocation(location);
	if err != nil {
		renderInternalServerError(w, err)
        return
	}
	if !find {
		hash := map[string]string{ "error": "Name not found"}
		renderClientError(w, http.StatusNotFound, hash)
		return
	}
	model.DeleteLocation(location)
    w.WriteHeader(http.StatusOK)
}

func HandleGetLocation(w http.ResponseWriter, r *http.Request) {
	location := html.EscapeString(r.URL.Path[len("/location/"):])
	location = strings.ToLower(location)
	log.Println("get:", location)
	find, err := model.FindLocation(location);
	if err != nil {
		renderInternalServerError(w, err)
        return
	}
	if !find {
		hash := map[string]string{ "error": "Name not found"}
		renderClientError(w, http.StatusNotFound, hash)
		return
	}
	weather, err := model.GetWeather(location)
	if err != nil {
		renderInternalServerError(w, err)
        return
	}

	log.Printf("request location: %v, weather: %v", location, weather)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
	w.Write([]byte(weather))
}

func LocationIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		HandleLocationList(w, r)
	} else if r.Method == "POST" {
		HandleCreateLocation(w, r)
	}
}

func LocationElement(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		HandleDeleteLocation(w, r)
	} else if r.Method == "GET" {
		HandleGetLocation(w, r)
	}
}

// internal
func renderInternalServerError(w http.ResponseWriter, err error) {
	log.Println("internal server error: ", err)
    hash := map[string]string{ "error": "internal server error"}
    render(w, http.StatusInternalServerError, hash)
}

func renderClientError(w http.ResponseWriter, code int, body interface{}) {
	log.Println("client error: ", body, code)
	render(w, code, body)
}

func renderOK(w http.ResponseWriter, code int, body interface{}) {
	render(w, code, body)
}

func render(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
        panic(err)
    }
}