package model

import (
	"LidGeneration/static"
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
)

type jsonResponseDictionary struct {
	Array []struct {
		Synonyms []struct {
			Text string `json:"text"`
		} `json:"tr"`
	} `json:"def"`
}

/*
Функция разбивает поисковую фразу на отдельные слова и анализирует текст title и description на эти слова и их синонимы.
Критерии строго анализа: все слова или их синонимы из фразы должны быть либо в title либо в description
Критерии нестрого анализа: хотя бы 1 слово или его синоним должно находится в title или description
*/
func textAnalyzerModel(keyPhrase, title, description string, strictSearch bool) (result bool) {
	var needle string
	// Приводим все к нижнему регистру
	keyPhrase = strings.ToLower(keyPhrase)
	title = strings.ToLower(title)
	description = strings.ToLower(description)
	// разбиваем фразу на отдельные слова и убираем предлоги
	keyWordsArr := removePretexts(keyPhrase)
	// получаем синонимы ключевых слов
	keyWordsArrWithSynonyms := getSynonyms(keyWordsArr)
	// поиск по ключевым словам
	for _, synonymsArr := range keyWordsArrWithSynonyms {
		if strictSearch == true {
			// если задан строгий анализ
			needle += "(" + strings.Join(synonymsArr, "|") + ")+.*" // хотя бы 1 раз должен быть каждый ключ либо его синоним
		} else {
			// если задан не строгий анализ
			needle += "(" + strings.Join(synonymsArr, "|") + ")|" // хотя бы 1 раз должен быть любой ключ или его синоним
		}
	}
	// Убираем последний символ |
	if strictSearch == false {
		needle = needle[:len(needle)-1]
	}
	// Ищем совпадения
	regExp, _ := regexp.Compile("(?i)" + needle)
	if regExp.MatchString(title) || regExp.MatchString(description) {
		result = true
	}

	return result
}

/*
Функция убирает предлоги из ключевой фразы и возвращает оставшиеся слова в виде массива
*/
func removePretexts(phrase string) (result []string) {
	var foundPretext bool
	var pretexts = [20]string{"без", "до", "из", "к", "в", "на", "по", "о", "от", "перед", "при", "через", "с", "у", "за", "над", "об", "под", "про", "для"}
	var words = strings.Split(phrase, " ")
	for _, word := range words {
		foundPretext = false
		// пропускаем предлоги
		for _, pretext := range pretexts {
			if word == pretext {
				foundPretext = true
				break
			}
		}
		if foundPretext == false {
			result = append(result, word)
		}
	}
	return result
}

/*
Функция запрашивает синонимы слов у словаря через api и возвращает массив синонимов с удаленными окончаниями
*/
func getSynonyms(words []string) map[string][]string {
	var result = map[string][]string{}
	queryParams := url.Values{}
	queryParams.Set("key", static.DICTIONARY_KEY)
	queryParams.Set("lang", static.DICTIONARY_LANG)
	// отправляем запрос
	for _, word := range words {
		// убираем окончание у ключа
		wordWithoutEnding := removeWordEnding(word)
		// записываем в массив
		var synonyms = []string{wordWithoutEnding[0]}
		// получаем синонимы
		queryParams.Set("text", word)
		resp, err := httpRequest(static.DICTIONARY_API+queryParams.Encode(), "GET", nil, "", false)
		if err == nil {
			jsonResp := new(jsonResponseDictionary)
			err = json.Unmarshal(resp, &jsonResp)
			if err == nil && len(jsonResp.Array) > 0 && len(jsonResp.Array[0].Synonyms) > 0 {
				for _, synonym := range jsonResp.Array[0].Synonyms {
					synonyms = append(synonyms, synonym.Text)
				}
				// убираем окончание у синонимов
				synonyms = removeWordEnding(synonyms...)
			}
		}
		// записываем результат
		result[word] = synonyms
	}
	return result
}

/*
Функция является упрощенной формой "Стеммера Портера". Она убирает окончание у слов и возвращает результат.
Это необходимо для анализа выдачи поисковой системы, поиска ключей в Title и Description.
*/
func removeWordEnding(words ...string) (result []string) {
	regExp, _ := regexp.Compile("[аоийеёэыуюя]")
	for _, word := range words {
		// Пздц с кириллицей в go
		word := []rune(word)
		// убираем окончания
		for i := len(word) - 1; i > 0; i-- {
			if regExp.MatchString(string(word[i])) {
				word = word[:len(word)-1]
			} else {
				break
			}
		}
		result = append(result, string(word))
	}
	return result
}

/*
Функция транслитерирует русские символы на английские
*/
func translitRuToEn(word string) (result string) {
	word = strings.ToLower(word)
	wordInRune := []rune(word)
	for i := 0; i < len(wordInRune); i++ {
		switch wordInRune[i] {
		case 'а':
			result += "a"
			continue
		case 'б':
			result += "b"
			continue
		case 'в':
			result += "v"
			continue
		case 'г':
			result += "g"
			continue
		case 'д':
			result += "d"
			continue
		case 'е':
			result += "e"
			continue
		case 'ё':
			result += "e"
			continue
		case 'ж':
			result += "zh"
			continue
		case 'з':
			result += "z"
			continue
		case 'и':
			result += "i"
			continue
		case 'й':
			result += "y"
			continue
		case 'к':
			result += "k"
			continue
		case 'л':
			result += "l"
			continue
		case 'м':
			result += "m"
			continue
		case 'н':
			result += "n"
			continue
		case 'о':
			result += "o"
			continue
		case 'п':
			result += "p"
			continue
		case 'р':
			result += "r"
			continue
		case 'с':
			result += "s"
			continue
		case 'т':
			result += "t"
			continue
		case 'у':
			result += "u"
			continue
		case 'ф':
			result += "f"
			continue
		case 'х':
			result += "h"
			continue
		case 'ц':
			result += "ts"
			continue
		case 'ч':
			result += "ch"
			continue
		case 'ш':
			result += "sh"
			continue
		case 'щ':
			result += "sch"
			continue
		case 'ъ':
			result += "''"
			continue
		case 'ы':
			result += "y"
			continue
		case 'ь':
			result += "''"
			continue
		case 'э':
			result += "e"
			continue
		case 'ю':
			result += "yu"
			continue
		case 'я':
			result += "ya"
			continue
		default:
			continue
		}
	}
	return result
}
