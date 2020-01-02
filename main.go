package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Map struct {
	Info    Info      `json:"info"`
	Options Options   `json:"options"`
	Results []Results `json:"results"`
}
type Info struct {
	Statuscode int       `json:"statuscode"`
	Copyright  Copyright `json:"copyright"`
	Messages   []string  `json:"messages"`
}
type Copyright struct {
	Text         string `json:"text"`
	ImageURL     string `json:"imageUrl"`
	ImageAltText string `json:"imageAltText"`
}
type Options struct {
	MaxResults        int  `json:"maxResults"`
	ThumbMaps         bool `json:"thumbMaps"`
	IgnoreLatLngInput bool `json:"ignoreLatLngInput"`
}
type Results struct {
	ProvidedLocation ProvidedLocation `json:"providedLocation"`
	Locations        []Location       `json:"locations"`
}
type ProvidedLocation struct {
	Location string `json:"location"`
}
type Location struct {
	Street             string       `json:"street"`
	AdminArea6         string       `json:"adminArea6"`
	AdminArea6Type     string       `json:"adminArea6Type"`
	AdminArea5         string       `json:"adminArea5"`
	AdminArea5Type     string       `json:"adminArea5Type"`
	AdminArea4         string       `json:"adminArea4"`
	AdminArea4Type     string       `json:"adminArea4Type"`
	AdminArea3         string       `json:"adminArea3"`
	AdminArea3Type     string       `json:"adminArea3Type"`
	AdminArea1         string       `json:"adminArea1"`
	AdminArea1Type     string       `json:"adminArea1Type"`
	PostalCode         string       `json:"postalCode"`
	GeocodeQualityCode string       `json:"geocodeQualityCode"`
	GeocodeQuality     string       `json:"geocodeQuality"`
	DragPoint          bool         `json:"dragPoint"`
	SideOfStreet       string       `json:"sideOfStreet"`
	LinkID             string       `json:"linkId"`
	UnknownInput       string       `json:"unknownInput"`
	Type               string       `json:"type"`
	LatLng             LatLng       `json:"latLng"`
	DisplayLatLn       DisplayLatLn `json:"displayLatLng"`
	MapURL             string       `json:"mapUrl"`
}
type DisplayLatLn struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func convert(location string) string {
	a := strings.Split(location, ", ")
	return strings.Join(a, "+")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/map", getmap)

	port := ":8000"

	fmt.Printf("Listen on port %s\n", port)
	//this specifyis that it should listen on port :8080, it will blick until the program is terminated
	//the ListenAndServe return a error that way the log.Fatal is there to output the error if it occurs
	log.Fatal(http.ListenAndServe(port, nil))
}

func setStruct(converted string, tr chan Map) {
	result := "http://open.mapquestapi.com/geocoding/v1/address?key=" + "6ZAjU5GDUjzJIaU1aMsMEf19AlG4hNCx&location=" + converted
	response, err := http.Get(result)
	if err != nil {
		log.Panic(err)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic(err)
	}
	save := Map{}
	err = json.Unmarshal(data, &save)
	if err != nil {
		log.Fatal(err)
	}
	tr <- save
}

var mapp Map

//http://open.mapquestapi.com/geocoding/v1/address=/R.+da+Mouraria+3,+9000-047+Funchal
// ipa for google maps :)
func handler(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := templates.Execute(w, ""); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	converted := convert(r.FormValue("str"))
	tr := make(chan Map)
	if converted != "" {
		go setStruct(converted, tr)

		mapp = <-tr
		fmt.Println(mapp)
	}
}

func getmap(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("map.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := templates.Execute(w, mapp.Results[0].Locations[0]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
