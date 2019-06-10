package model

import (
	"LidGeneration/static"
	"bufio"
	"net/http"
	"os"
	"strings"
	"text/template"
)

/*
Функция забирает все html файлы в папке view
*/
func ViewModel() *template.Template {
	viewFolder, _ := os.Open(static.VIEW_FOLDER_PATH)
	defer viewFolder.Close()
	viewPathRaw, _ := viewFolder.Readdir(-1)

	var viewPaths []string
	for _, pathInfo := range viewPathRaw {
		if !pathInfo.IsDir() {
			viewPaths = append(viewPaths, static.VIEW_FOLDER_PATH+"/"+pathInfo.Name())
		}
	}

	// Подключаем css и js
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)

	result, _ := template.ParseFiles(viewPaths...)

	return result
}

/*
Функция задает content-type в header для подключаемых скриптов
*/
func serveResource(w http.ResponseWriter, req *http.Request) {
	// папка, где лежат стили
	path := static.VIEW_FOLDER_PATH + req.URL.Path
	// задаем content-type для файлов
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".js") {
		contentType = "text/javascript"
	} else {
		contentType = "text/plain"
	}
	// Добавляем заголовок
	f, err := os.Open(path)
	if err == nil {
		defer f.Close()
		w.Header().Add("Content-Type", contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
