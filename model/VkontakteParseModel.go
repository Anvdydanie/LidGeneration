package model

import (
	"LidGeneration/static"
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type vkJsonResponse struct {
	Response struct {
		Counts int `json:"count"`
		Groups []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"items"`
	} `json:"response"`
	Error string `json:"error"`
}

/*
Функция парсит через api ВКонтакте список групп по ключевым словам
*/
func VkontakteParseModel(theme, city string, strictSearch bool, arrCh chan map[string][]map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	var result map[string][]map[string]string
	var resultFiltered, resultFull []map[string]string
	var cityId string
	// получаем id города
	citiesId, err := httpRequest(static.VK_API_GET_CITIES_ID+city, "GET", nil, "", false)
	if err == nil {
		// сделаем небольшую правку для соответствия json формату
		citiesId = []byte(strings.ReplaceAll(string(citiesId), "'", "\""))
		// разбираем ответ
		var vkCitiesArr [][]string
		// [[ID, Название города, Область], ...]
		err = json.Unmarshal(citiesId, &vkCitiesArr)
		if err == nil {
			cityId = vkCitiesArr[0][0]
		} else {
			err = errors.New("При получении ID города " + city + " в ВКонтакте была получена ошибка " + err.Error())
		}
	}
	// Получаем список подходящих по теме групп
	if err == nil {
		// убираем предлоги, разбиваем фразу на слова и удаляем окончания
		searchWords := removePretexts(theme)
		searchWords = removeWordEnding(searchWords...)
		// составляем запрос
		var params = url.Values{}
		params.Set("type", static.VK_API_GROUP_TYPE)
		params.Set("city_id", cityId)
		params.Set("count", static.VK_API_GROUPS_COUNT)
		params.Set("access_token", static.VK_API_TOKEN)
		params.Set("v", static.VK_API_VERSION)
		// Ищем группы по каждому ключевому слову, т.к. по фразе они редко находятся
		for _, searchWord := range searchWords {
			params.Set("q", searchWord)
			vkGroupsList, err := httpRequest(static.VK_API_GET_GROUPS+params.Encode(), "GET", nil, "", false)
			if err == nil {
				var vkGroupsArr = new(vkJsonResponse)
				err = json.Unmarshal(vkGroupsList, &vkGroupsArr)
				if err == nil {
					if vkGroupsArr.Response.Counts > 0 {
						for _, group := range vkGroupsArr.Response.Groups {
							// записываем все результаты
							resultFull = append(resultFull, map[string]string{
								"title":   group.Name,
								"domain":  "vk.com/" + strconv.Itoa(group.Id),
								"fullUrl": static.VK_GROUP_URL + strconv.Itoa(group.Id),
							})
							// записываем только релевантные результаты
							var resultIsRelevant = textAnalyzerModel(theme, group.Name, "", strictSearch)
							if resultIsRelevant == true {
								resultFiltered = append(resultFiltered, map[string]string{
									"title":   group.Name,
									"domain":  "vk.com/" + strconv.Itoa(group.Id),
									"fullUrl": static.VK_GROUP_URL + strconv.Itoa(group.Id),
								})
							}
						}
					}
				} else {
					err = errors.New("При получении списка групп в ВКонтакте была получена ошибка: " + err.Error())
				}
			}
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
