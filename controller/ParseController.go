package controller

import (
	"LidGeneration/model"
	"LidGeneration/static"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

type jsonFromFrontend struct {
	Theme        string
	CityName     string
	StrictSearch bool
}

/*
Функция возвращает общий массив выдачи всех сервисов.
*/
func ParseController(w http.ResponseWriter, req *http.Request) {
	/* фильтр агрегаторов и ненужных сайтов */
	var filteredUrls = map[string]bool{
		"avito.ru":      true,
		"Яндекс.Карты":  true,
		"Яндекс.Маркет": true,
		"yandex.ru":     true,
		"2gis.ru":       true,
	}

	if req.Method == "POST" {
		contents, _ := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		// Разбираем json полученный при ajax запросе на парсинг
		var reqData = new(jsonFromFrontend)
		err := json.Unmarshal(contents, &reqData)
		if err == nil {
			chanCount := 5
			var yaSearchChan = make(chan map[string][]map[string]string)
			var yaAdvertChan = make(chan map[string][]map[string]string)
			var googleChan = make(chan map[string][]map[string]string)
			var vkChan = make(chan map[string][]map[string]string)
			var avitoChan = make(chan map[string][]map[string]string)
			var wg sync.WaitGroup
			wg.Add(chanCount)
			// Парсинг поисковой выдачи яндекса
			go model.YandexXmlParseModel(reqData.Theme, reqData.CityName, filteredUrls, reqData.StrictSearch, yaSearchChan, &wg)
			// Парсинг рекламной выдачи яндекса
			go model.YandexSearchParseModel(reqData.Theme, reqData.CityName, filteredUrls, reqData.StrictSearch, static.ADVERT_TYPE, yaAdvertChan, &wg)
			// Парсинг данных по googleapis. Легальный способ парсинга выдачи Гугла, но с ограничением в 100 запров в сутки
			go model.GoogleJsonParseModel(reqData.Theme, reqData.CityName, filteredUrls, reqData.StrictSearch, googleChan, &wg)
			// поиск групп в VK
			go model.VkontakteParseModel(reqData.Theme, reqData.CityName, reqData.StrictSearch, vkChan, &wg)
			// поиск объявлений авито
			go model.AvitoParseModel(reqData.Theme, reqData.CityName, reqData.StrictSearch, avitoChan, &wg)
			// получаем результат
			_result := map[string]map[string][]map[string]string{
				"yandexSearch": <-yaSearchChan,
				"googleSearch": <-googleChan,
				"yandexAdvert": <-yaAdvertChan,
				"vkGroups":     <-vkChan,
				"avito":        <-avitoChan,
			}
			wg.Wait()

			result, err := json.Marshal(_result)
			if err == nil {
				w.Write(result)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
