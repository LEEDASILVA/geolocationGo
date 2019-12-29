package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Location struct {
	Street             string
	AdminArea6         string
	AdminArea6Type     string
	AdminArea5         string
	AdminArea5Type     string
	AdminArea4         string
	AdminArea4Type     string
	AdminArea3         string
	AdminArea3Type     string
	AdminArea1         string
	AdminArea1Type     string
	PostalCode         string
	GeocodeQualityCode string
	GeocodeQuality     string
	DragPoint          bool
	SideOfStreet       string
	LinkID             string
	UnknownInput       string
	Type               string
	LatLng             LatLng
	DisplayLatLn       DisplayLatLn
	MapURL             string
}
type DisplayLatLn struct {
	Lat int
	Lng int
}
type LatLng struct {
	Lat int
	Lng int
}
type Info struct {
	Statuscode int
	Copyright  Copyright
	Messages   []string
}
type Copyright struct {
	Text         string
	ImageURL     string
	ImageAltText string
}
type Options struct {
	MaxResults        int
	ThumbMaps         bool
	IgnoreLatLngInput bool
}
type ProvidedLocation struct {
	Location string
}
type Results struct {
	ProvidedLocation ProvidedLocation
	Locations        []Location
}
type Map struct {
	Info    Info
	Options Options
	Results []Results
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

var save Map

//http://open.mapquestapi.com/geocoding/v1/address=/R.+da+Mouraria+3,+9000-047+Funchal
// ipa for google maps :)
func handler(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := templates.Execute(w, save); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	converted := convert(r.FormValue("str"))
	if converted != "" {
		result := "http://open.mapquestapi.com/geocoding/v1/address?key=" + "6ZAjU5GDUjzJIaU1aMsMEf19AlG4hNCx&location=" + converted
		fmt.Println(result)
		response, err := http.Get(result)
		if err != nil {
			log.Panic(err)
		}
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Panic(err)
		}

		json.Unmarshal(data, &save)
		fmt.Println(save)
	}
}

func getmap(w http.ResponseWriter, r *http.Request) {
	templates, err := template.ParseFiles("map.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := templates.Execute(w, save); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
