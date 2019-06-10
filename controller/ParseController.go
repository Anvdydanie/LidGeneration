package controller

import (
	"LidGeneration/model"
	"LidGeneration/static"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
			// TODO горутины
			// Парсинг данных по яндекс.xml. Яндекс предоставляет легальный способ парсинга выдачи, но с ограничениями на количество запров в сутки
			yaSitesList, _ := model.YandexXmlParseModel(reqData.Theme, reqData.CityName, filteredUrls, reqData.StrictSearch)
			if yaSitesList == nil {
				// нелегально парсим данные по яндексу, в случае, если лимит в яндекс.xml исчерпан или сервис недоступен
				yaSitesList, _ = model.YandexSearchParseModel(reqData.Theme, reqData.CityName, filteredUrls, reqData.StrictSearch, static.SEARCH_TYPE)
			}
			// Парсинг рекламной выдачи яндекса
			yaAdvList, _ := model.YandexSearchParseModel(reqData.Theme, reqData.CityName, filteredUrls, reqData.StrictSearch, static.ADVERT_TYPE)
			// Парсинг данных по googleapis. Легальный способ парсинга выдачи Гугла, но с ограничением в 100 запров в сутки
			goSitesList, _ := model.GoogleJsonParseModel(reqData.Theme, reqData.CityName, filteredUrls, reqData.StrictSearch)
			// поиск групп в VK
			vkGroupList, _ := model.VkontakteParseModel(reqData.Theme, reqData.CityName, reqData.StrictSearch)
			// поиск объявлений авито
			avitoList, _ := model.AvitoParseModel(reqData.Theme, reqData.CityName, reqData.StrictSearch)

			// объединяем все полученные результаты в 1 массив
			_result := map[string]map[string][]map[string]string{
				"yandexSearch": yaSitesList,
				"googleSearch": goSitesList,
				"yandexAdvert": yaAdvList,
				"vkGroups":     vkGroupList,
				"avito":        avitoList,
			}
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
