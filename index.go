package main

import (
	"LidGeneration/controller"
	"fmt"
	"net/http"
)

/* Критерии работы алгоритма:
1. Выборка результатов:
 - должна быть строгая? Пример запроса: пошив штор. Выборка только с фразами где присутствует: пошив, на заказ, под ключ и др.
 - или должна показывать все результаты со словом шторы: пошив, салон, покупка, купить и т.д.?
2. Дать возможность исключать url.
3. Ограничение парсинга в первые 50 результатов поиска
4. Возможность выборки:
 4.1 Строгая: в результатах либо в title либо в description обязательно должны присутствовать оба слова в разных падежах из темы.
 4.2 Мягкая: хотя бы 1 слово из темы должны быть либо в title либо в description
 *В будущем будут учитываться синонимы.
5.
*/

func main() {
	// Главная страница
	http.HandleFunc("/", controller.ViewController)

	// Запрос списка городов для автокомплита
	http.HandleFunc("/getCitiesList", controller.CitiesController)

	// Запрос парсинга выдачи поисковиков и агрегаторов по теме
	http.HandleFunc("/parseSearchEngines", controller.ParseController)

	// запускаем вебсервер
	err := http.ListenAndServe(":9000", nil)
	if err == nil {
		fmt.Println("Server is listening")
	} else {
		fmt.Println(err)
	}
}
