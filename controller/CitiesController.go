package controller

import (
	"LidGeneration/model"
	"io/ioutil"
	"net/http"
)

func CitiesController(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		defer req.Body.Close()
		contents, err := ioutil.ReadAll(req.Body)
		if err == nil {
			citiesList, err := model.CitiesModel(contents)
			if err == nil {
				w.Write(citiesList)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
