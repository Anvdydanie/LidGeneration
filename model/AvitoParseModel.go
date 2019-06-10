package model

import (
	"LidGeneration/static"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

/*
Функция парсит список предложений авито по поисковой фразе
*/
func AvitoParseModel(theme, city string, strictSearch bool) (result map[string][]map[string]string, err error) {
	var resultFiltered, resultFull []map[string]string
	// получаем транслит названия города для поиска
	var cityInTranslit = translitRuToEn(city)
	// составляем запрос
	var params = url.Values{}
	params.Set("s_trg", static.AVITO_CATEGORY_ID)
	params.Set("q", theme)
	requestUrl := static.AVITO_SEARCH_URL + cityInTranslit + "?" + params.Encode()
	// Ищем объявления в авито
	resp, err := httpRequest(requestUrl, "GET", nil, "", true)
	if err == nil {
		// Приводим результат парсинга в формат IoReader для goquery
		htmlIoRead := strings.NewReader(string(resp))
		// обращаемся к элементам html как в jquery
		doc, err := goquery.NewDocumentFromReader(htmlIoRead)
		if err == nil {
			//.item.item_table - класс li в выдаче
			doc.Find(".item.item_table").Each(func(i int, s *goquery.Selection) {
				// получаем ссылку на главную страницу сайта из выдачи
				itemUrl, _ := s.Find(".item-description-title-link").Attr("href")
				title, _ := s.Find(".item-description-title-link").Attr("title")
				description := s.Find(".data").Text()
				price := s.Find(".about .price").Text()
				// полный ответ, содержит превые 50 результатов
				resultFull = append(resultFull, map[string]string{
					"title":       title,
					"domain":      "avito.ru/" + itemUrl,
					"fullUrl":     static.AVITO_SEARCH_URL + itemUrl,
					"description": description + price,
				})
				// Анализируем ответ на наличие ключевых слов с использованием синонимов
				var resultIsRelevant = textAnalyzerModel(theme, title, description, strictSearch)
				if resultIsRelevant == true {
					resultFiltered = append(resultFiltered, map[string]string{
						"title":       title,
						"domain":      "avito.ru/" + itemUrl,
						"fullUrl":     static.AVITO_SEARCH_URL + itemUrl,
						"description": description + price,
					})
				}
			})
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
