package model

import (
	"LidGeneration/static"
	"encoding/xml"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type yandexXmlResponse struct {
	Response struct {
		Group []struct {
			Title       string `xml:"title"`
			Url         string `xml:"url"`
			Domain      string `xml:"domain"`
			Description string `xml:"passages>passage"`
		} `xml:"group>doc"`
	} `xml:"response>results>grouping"`
}

/*
Функция парсит результаты выдачи яндекса через их XML сервис
*/
func YandexXmlParseModel(theme, city string, filteredUrls map[string]bool, strictSearch bool, arrCh chan map[string][]map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	var result map[string][]map[string]string
	var err error
	// Результаты полной и фильтрованной выборки
	var resultFiltered, resultFull []map[string]string
	// Параметры запроса
	queryParams := url.Values{}
	queryParams.Add("user", static.YANDEX_SEARCH_API_XML_USER)
	queryParams.Add("key", static.YANDEX_SEARCH_API_XML_KEY)
	queryParams.Add("query", theme+" "+city)
	queryParams.Add("sortby", static.YANDEX_SEARCH_API_XML_SORT)
	queryParams.Add("filter", static.YANDEX_SEARCH_API_XML_FILTER)
	queryParams.Add("groupby", static.YANDEX_SEARCH_API_XML_GROUP)
	// Отправляем запрос
	resp, err := httpRequest(static.YANDEX_SEARCH_API_XML_URL+queryParams.Encode(), "GET", nil, "", false)
	// Яндекс впихнул в ответ, тег hlword с выделением ключей в title и description, в итоге распарсить xml нормально не получается
	// Приходится чистить варварским способом
	resp = []byte(strings.ReplaceAll(string(resp), "<hlword>", ""))
	resp = []byte(strings.ReplaceAll(string(resp), "</hlword>", ""))
	// Разбираем xml
	if err == nil {
		var yaXml = new(yandexXmlResponse)
		if err = xml.Unmarshal(resp, &yaXml); err == nil {
			for key, group := range yaXml.Response.Group {
				// полный ответ, содержит все 50 результатов
				resultFull = append(resultFull, map[string]string{
					"domain":      group.Domain,
					"fullUrl":     group.Url,
					"title":       group.Title,
					"description": group.Description,
					"position":    strconv.Itoa(key + 1),
				})
				// фильтрованный ответ, убирает агрегаторы из выдачи и повторы
				if filteredUrls[group.Domain] == false {
					// Анализируем ответ на наличие ключевых слов с использованием синонимов
					var resultIsRelevant = textAnalyzerModel(theme, group.Title, group.Description, strictSearch)
					if resultIsRelevant == true {
						resultFiltered = append(resultFiltered, map[string]string{
							"domain":      group.Domain,
							"fullUrl":     group.Url,
							"title":       group.Title,
							"description": group.Description,
							"position":    strconv.Itoa(key + 1),
						})
					}
				}
			}
		} else {
			err = errors.New("В ходе парсинга XML выдачи яндекса, был получен неожиданный ответ: " + err.Error())
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
