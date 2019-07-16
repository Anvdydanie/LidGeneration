package model

import (
	"LidGeneration/static"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"sync"
)

type jsonResponseGoogle struct {
	Items []struct {
		Title  string `json:"title"`
		Descr  string `json:"snippet"`
		FLink  string `json:"link"`
		Domain string `json:"displayLink"`
	} `json:"items"`
}

/*
Функция легально парсит выдачу Google c помощью api сервиса googleapis.com. Ограничение в 100 запросов в сутки
*/
func GoogleJsonParseModel(theme, city string, filteredUrls map[string]bool, strictSearch bool, arrCh chan map[string][]map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	var result map[string][]map[string]string
	// Результаты полной и фильтрованной выборки
	var resultFiltered, resultFull []map[string]string
	// Параметры запроса
	queryParams := url.Values{}
	queryParams.Set("key", static.GOOGLE_API_JSON_KEY)
	queryParams.Set("cx", static.GOOGLE_API_JSON_SID)
	queryParams.Set("q", theme+" "+city)
	// Гугл api выдает максимум 10 результатов, поэтому приходится выбирать с какого результата делать выборку
	for i := 0; i < static.MAX_PAGES_TO_PARSE; i++ {
		queryParams.Set("start", strconv.Itoa(i*10+1))
		urlRequest := static.GOOGLE_SEARCH_API_JSON_URL + queryParams.Encode()
		// Отправляем запрос
		resp, err := httpRequest(urlRequest, "GET", nil, "", false)
		if err == nil {
			jsonResp := new(jsonResponseGoogle)
			err := json.Unmarshal(resp, &jsonResp)
			if err == nil {
				for key, value := range jsonResp.Items {
					// полный ответ, содержит все 50 результатов
					resultFull = append(resultFull, map[string]string{
						"domain":      value.Domain,
						"fullUrl":     value.FLink,
						"title":       value.Title,
						"description": value.Descr,
						"position":    strconv.Itoa(i*10 + key + 1),
					})
					// фильтрованный ответ, убирает агрегаторы из выдачи
					if filteredUrls[value.Domain] == false {
						// Анализируем текст
						var resultIsRelevant = textAnalyzerModel(theme, value.Title, value.Descr, strictSearch)
						if resultIsRelevant == true {
							resultFiltered = append(resultFiltered, map[string]string{
								"domain":      value.Domain,
								"fullUrl":     value.FLink,
								"title":       value.Title,
								"description": value.Descr,
								"position":    strconv.Itoa(i*10 + key + 1),
							})
						}
					}
				}
			} else {
				// TODO лимит закончился, меняем api key или возвращаем ошибку
				err = errors.New("Лимит googleapis.com закончился: " + err.Error())
			}
		} else {
			err = errors.New("Не удалось сделать запрос к googleapis.com: " + err.Error())
		}
	}
	// Совмещаем полную и фильтрованную выдачу
	if len(resultFull) > 0 {
		result = map[string][]map[string]string{
			"filteredResponse": resultFiltered,
			"fullResponse":     resultFull,
		}
	}
	arrCh <- result
}
