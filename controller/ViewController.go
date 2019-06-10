package controller

import (
	"LidGeneration/model"
	"net/http"
)

// получаем все шаблоны
var templates = model.ViewModel()

func ViewController(w http.ResponseWriter, req *http.Request) {
	// исключение для главной страницы
	requestedFile := req.URL.Path[1:] // 1 символ ( / ) не нужен
	if requestedFile == "" {
		requestedFile = "main"
	}
	// отдаем шаблон при наличии или возвращаем 404
	view := templates.Lookup(requestedFile + ".html")
	if view != nil {
		view.Execute(w, nil)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
