package model

import (
	"LidGeneration/static"
	"encoding/xml"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strconv"
	"strings"
)

// Определяем интерфейс YandexSearch для поиска в поисковой и рекламной выдаче яндекса
type YandexSearch interface {
	yaSearch() (result map[string][]map[string]string, err error)
}

// Объявляем метод для реализации
func parseSearch(searchType YandexSearch) (result map[string][]map[string]string, err error) {
	result, err = searchType.yaSearch()
	return result, err
}

// Создаем структуру под поиск в рекламных объявлениях и поисковой выдаче
type yaSearchParams struct {
	Url          string
	StrictSearch bool
	staticSearchParams
}
type staticSearchParams struct {
	City         string
	Theme        string
	FilteredUrls map[string]bool
	PagesToParse int
}

// Описываем метод структуры yaSearchParams
func (sType yaSearchParams) yaSearch() (result map[string][]map[string]string, err error) {
	var city = sType.staticSearchParams.City
	var theme = sType.staticSearchParams.Theme
	var filteredUrls = sType.staticSearchParams.FilteredUrls
	var pagesToParse = sType.staticSearchParams.PagesToParse
	var strictSearch = sType.StrictSearch
	var searchUrl = sType.Url
	// Результаты полной и фильтрованной выборки
	var resultFiltered, resultFull []map[string]string
	// Получаем id города через яндекс.погоду
	cityId, err := getYandexCitiesId(city)
	// Если город найден, парсим выдачу яндекса
	if err == nil && cityId != "" {
		queryParams := url.Values{}
		queryParams.Set("text", theme+" "+city)
		queryParams.Set("lr", cityId)
		// Постранично запрашиваем результаты
		for i := 0; i < pagesToParse; i++ {
			// добавляем страницы в параметрах
			if i > 0 {
				queryParams.Set("p", strconv.Itoa(i))
			}
			// Запрос выдачи яндекса
			yandexParse, err := httpRequest(searchUrl+queryParams.Encode(), "GET", nil, "", true)
			if err == nil {
				// Приводим результат парсинга в формат IoReader для goquery
				htmlIoRead := strings.NewReader(string(yandexParse))
				// обращаемся к элементам html как в jquery
				doc, err := goquery.NewDocumentFromReader(htmlIoRead)
				if err == nil {
					// останавливаем цикл в случае, если количество результатов было достигнуто раньше окончания цикла
					if doc.Find(".main__content .misspell__message").Text() != "" {
						break
					}
					//.serp-item - класс li в выдаче яндекса
					doc.Find(".serp-item").Each(func(j int, s *goquery.Selection) {
						// получаем ссылку на главную страницу сайта из выдачи
						domain := s.Find(".link.link_theme_outer.path__item.i-bem").Children().Text()
						title := s.Find(".organic__url-text").Text()
						description := s.Find(".organic__content-wrapper").Text()
						// полный ответ, содержит все 50 результатов
						resultFull = append(resultFull, map[string]string{
							"domain":      domain,
							"fullUrl":     domain,
							"title":       title,
							"description": description,
							"position":    strconv.Itoa(i*10 + j + 1),
						})
						// фильтрованный ответ, убирает агрегаторы из выдачи и повторы
						if filteredUrls[domain] == false {
							// Анализируем ответ на наличие ключевых слов с использованием синонимов
							var resultIsRelevant = textAnalyzerModel(theme, title, description, strictSearch)
							if resultIsRelevant == true {
								resultFiltered = append(resultFiltered, map[string]string{
									"domain":      domain,
									"fullUrl":     domain,
									"title":       title,
									"description": description,
									"position":    strconv.Itoa(i*10 + j + 1),
								})
							}
						}
					})
					// останавливаем цикл в случае, если количество результатов было достигнуто раньше окончания цикла
					// это возможно, если в настройках яндекса указан показ результатов выдачи более 10 на страницу
					if len(resultFull) > (pagesToParse-1)*10 {
						break
					}
				} else {
					err = errors.New("Проблема с загрузкой DOM выдачи яндекса в goquery: " + err.Error())
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
	return result, err
}

// Структура для разбора xml
type citiesStruct struct {
	Countries []struct {
		Name   string `xml:"name,attr"`
		Cities []struct {
			CityName string `xml:",chardata"`
			CityId   string `xml:"region,attr"`
		} `xml:"city"`
	} `xml:"country"`
}

/*
Функция парсит топ-50 выдачу яндекса, анализирует ее по ключевым словам. Возвращает полный и релевантный ключам результат
*/
func YandexSearchParseModel(theme, city string, filteredUrls map[string]bool, strictSearch bool, typeOfSearch string) (result map[string][]map[string]string, err error) {
	// Статичные параметры для любых запросов
	var staticParams = staticSearchParams{
		city,
		theme,
		filteredUrls,
		static.MAX_PAGES_TO_PARSE,
	}
	var searchParams yaSearchParams
	if typeOfSearch == static.SEARCH_TYPE {
		// запрос для поиска яндекса
		searchParams = yaSearchParams{static.YANDEX_SEARCH_URL, strictSearch, staticParams}
	} else if typeOfSearch == static.ADVERT_TYPE {
		// запрос для рекламных объявлений яндекса
		searchParams = yaSearchParams{static.YANDEX_ADVERT_URL, false, staticParams}
	} else {
		err = errors.New("Не распознан typeOfSearch. Необходимо использовать константу ..._TYPE из constants.go. ")
		return nil, err
	}
	// запрашиваем парсинг
	result, err = parseSearch(searchParams)

	return result, err
}

/*
Функция получает название города и возвращает его ID в системе Яндекса
*/
func getYandexCitiesId(city string) (result string, err error) {
	// Получаем id городов yandex
	yandexCities, err := httpRequest(static.YANDEX_CITIES_ID_URL, "GET", nil, "", false)
	if err == nil {
		var citiesData citiesStruct
		// парсим xml
		if err = xml.Unmarshal(yandexCities, &citiesData); err == nil {
			// Получаем id города
			for _, countryData := range citiesData.Countries {
				if countryData.Name == static.SEARCH_COUNTRY_NAME {
					for _, cityData := range countryData.Cities {
						if cityData.CityName == city {
							result = cityData.CityId
						}
					}
					break
				}
			}
		} else {
			err = errors.New("В ходе парсинга XML выдачи названий и ID городов яндекса, был получен неожиданный XML: " + err.Error())
		}
	}
	return result, err
}
