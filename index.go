package main

import (
	"LidGeneration/controller"
	"fmt"
	"net/http"
)

func main() {
	// Главная страница
	http.HandleFunc("/", controller.ViewController)

	// Запрос списка городов для автокомплита
	http.HandleFunc("/getCitiesList", controller.CitiesController)

	// Запрос парсинга выдачи поисковиков и агрегаторов по теме
	http.HandleFunc("/parseSearchEngines", controller.ParseController)

	// запускаем вебсервер
	err := http.ListenAndServe(":9000", nil)
	if err == nil {
		fmt.Println("Server is listening")
	} else {
		fmt.Println(err)
	}
}
