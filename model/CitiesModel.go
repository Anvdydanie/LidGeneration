package model

import (
	"LidGeneration/static"
	"encoding/json"
	"net/url"
)

type citiesJsonResp struct {
	Result []struct {
		Name string `json:"name"`
	} `json:"result"`
}

/*
Функция возвращает список городов России в json формате для автокомплита
*/
func CitiesModel(someLetters []byte) (result []byte, err error) {
	// параметры запроса
	queryParams := url.Values{}
	queryParams.Set("query", string(someLetters))
	queryParams.Set("contentType", static.CITIES_NAME_API_SETTLE_TYPE)
	queryParams.Set("typeCode", static.CITIES_NAME_API_TYPE_CODE)

	requestUrl := static.CITIES_NAME_API_URL + queryParams.Encode()
	resp, err := httpRequest(requestUrl, "GET", nil, "", false)
	if err == nil {
		cityStruct := new(citiesJsonResp)
		err = json.Unmarshal(resp, &cityStruct)
		if err == nil && len(cityStruct.Result) > 0 {
			var citiesName []string
			for _, cityData := range cityStruct.Result {
				citiesName = append(citiesName, cityData.Name)
			}
			result, err = json.Marshal(citiesName)
		}
	}
	return result, err
}
