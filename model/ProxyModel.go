package model

import (
	"LidGeneration/static"
	"database/sql"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type proxyData struct {
	id        int
	ip        string
	port      string
	proxyType string
	attempts  int
}

// TODO вынести в отдельные функции работу с БД

/*
Функция парсит списки прокси с разных сайтов, проверяет их на работоспособоности и сохраняет живые в базу данных (таблица: proxies)
*/
func ProxyModel() (err error) {
	start := time.Now()

	var sitesWithProxies = map[string]string{
		"online-proxy.ru": "http://online-proxy.ru/index.html",
		"spys.one":        "http://spys.one/free-proxy-list/RU/",
		/*	Еще сайты с прокси
			https://hidemyna.me/ru/proxy-list/?country=RU#list
			https://free.proxy-sale.com/
		*/
	}

	// соединяемся с БД
	db, err := sql.Open("sqlite3", static.DATABASE_FILE_PATH)
	defer db.Close()
	if err == nil {
		// Параллельно парсим сайты со списком прокси и проверяем их
		arrCh := make(chan []map[string]string, len(sitesWithProxies))
		var wg sync.WaitGroup
		wg.Add(len(sitesWithProxies))
		// получаем список прокси
		for domain, pageUrl := range sitesWithProxies {
			go parseProxiesByDomain(domain, pageUrl, arrCh, &wg)
		}
		for i := 0; i < len(sitesWithProxies); i++ {
			proxyList := <-arrCh
			if len(proxyList) > 0 {
				for _, value := range proxyList {
					if err = saveProxyInDb(db, value["ip"], value["port"], value["type"]); err != nil {
						err = errors.New("Запись в БД не осуществлена! " + err.Error())
					}
				}
			}
		}
		wg.Wait()
	}

	elapsed := time.Since(start)
	fmt.Println("Парсинг проксей завершен за", elapsed)

	return err
}

/*
Функция содержит парсер для каждого сайта из списка
*/
func parseProxiesByDomain(domain, pageUrl string, arrCh chan []map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	var result []map[string]string
	var proxiesAdded int
	var mainError error

	switch domain {
	case "online-proxy.ru":
		{
			respHtml, err := httpRequest(pageUrl, "GET", nil, "", false)
			if err == nil {
				fmt.Println("Парсим прокси с online-proxy.ru")
				htmlIoRead := strings.NewReader(string(respHtml))
				// обращаемся к элементам html как в jquery
				doc, _ := goquery.NewDocumentFromReader(htmlIoRead)
				// Удаляем ненужный фрагмент
				doc.Find(".content table:first-child").Remove()
				// разбираем строки
				doc.Find(".content table tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
					proxyIp := s.Find("td").Eq(1).Text()
					if proxyIp != "" { // Ограничиваем проверку первыми 100 проксями (там их тысячи)
						proxyPort := s.Find("td").Eq(2).Text()
						// Тип соединения с прокси
						proxyType := s.Find("td").Eq(3).Text()
						// проверяем рабочий ли прокси
						if proxyChecked(proxyIp, proxyPort, proxyType) == true {
							// Сохраняем живые прокси
							result = append(result, map[string]string{
								"ip":   proxyIp,
								"port": proxyPort,
								"type": proxyType,
							})
							// подсчитываем сколько живых прокси предоставил сайт
							proxiesAdded++
						}
					}
					// останавливаем цикл
					if i == 50 {
						return false
					}
					return true
				})
			} else {
				mainError = errors.New("Ошибка при парсинге online-proxy.ru: " + err.Error())
			}
			break
		}

	case "spys.one":
		{
			respHtml, err := httpRequest(pageUrl, "POST", nil, "", false)
			if err == nil {
				fmt.Println("Парсим прокси с spys.one")
				htmlIoRead := strings.NewReader(string(respHtml))
				// обращаемся к элементам html как в jquery
				doc, _ := goquery.NewDocumentFromReader(htmlIoRead)
				doc.Find("table:nth-child(3) table:nth-child(1) tr:nth-child(3)").NextAll().Each(func(i int, s *goquery.Selection) {
					// Текст js скрипта написан прямо в ip, поэтому мы его удаляем
					s.Find(`script[type="text/javascript"]`).Remove()
					proxyIp := s.Find("td").Eq(0).Find(".spy14").Text()
					// Порт по умолчанию
					proxyPort := "8080"
					// Тип соединения с прокси
					proxyType := s.Find("td").Eq(1).Find(".spy1").Text()
					// Проверяем прокси на работоспособность
					if proxyChecked(proxyIp, proxyPort, proxyType) == true {
						// Сохраняем живые прокси
						result = append(result, map[string]string{
							"ip":   proxyIp,
							"port": proxyPort,
							"type": proxyType,
						})
						// подсчитываем сколько живых прокси предоставил сайт
						proxiesAdded++
					}
				})
			} else {
				mainError = errors.New("Ошибка при парсинге spys.one: " + err.Error())
			}
			break
		}

	default:
		{
			// в случае, если ip и порты написаны через :
			/*regExp, _ := regexp.Compile("([0-9]{1,3}[.]){3}[0-9]{1,3}(<br>)?:[1-9]{2-6}")*/
			break
		}
	}

	if mainError != nil {
		fmt.Println(domain, mainError)
	} else {
		fmt.Println("Получено живых проксей с сайта "+domain+": ", proxiesAdded)
	}

	arrCh <- result
}

/*
Функция отправляет запрос к стороннему ресурсу для проверки состояния прокси
*/
func proxyChecked(ip, port, proxyType string) (result bool) {
	//fmt.Println("Проверяем прокси " + ip+":"+port + " на живучесть")
	const URL_TO_CHECK = "https://google.com"
	// Настраиваем клиент
	proxyUrl, err := url.Parse(proxyType + "://" + ip + ":" + port)
	var client *http.Client
	if err == nil {
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
			Timeout: time.Duration(5 * time.Second),
		}
		// Отправляем запрос
		resp, err := client.Get(URL_TO_CHECK)
		if err == nil {
			defer resp.Body.Close()
			//fmt.Println("Прокси живой:", ip, port)
			result = true
		} else {
			//fmt.Println("Прокси дохлый: " + err.Error())
		}
	} else {
		//fmt.Println("Возникла проблема с адресом прокси: " + err.Error())
	}
	return result
}

/*
Функция записывает в БД таблицу proxies живые прокси. Записываемые данные: ip, port, type (string)
*/
func saveProxyInDb(db *sql.DB, proxyIp, proxyPort, proxyType string) (err error) {
	var id int
	// ищем совпадения
	findRowErr := db.QueryRow("SELECT id FROM proxies WHERE ip = $1", proxyIp).Scan(&id)
	// если совпадений нет пишем в бд
	if findRowErr != nil || id == 0 {
		proxyType = strings.ToLower(proxyType)
		if _, err = db.Exec("INSERT INTO proxies (ip, port, `type`) VALUES ($1, $2, $3)", proxyIp, proxyPort, proxyType); err != nil {
			err = errors.New("Ошибка при записи прокси в БД: " + err.Error())
		}
	}
	return err
}

/*
Функция получает прокси из базы данных, проверяет его и возвращает. При 3 неудачных попытках (не подряд) использовать прокси, он удаляется из таблицы.
*/
func getProxyFromTable() (result map[string]string, err error) {
	db, err := sql.Open("sqlite3", static.DATABASE_FILE_PATH)
	defer db.Close()
	if err != nil {
		err = errors.New("Не удается получить доступ к Базе Данных: " + err.Error())
	}

	rows, err := db.Query("SELECT id, ip, port, `type`, attempts FROM proxies")
	if err != nil {
		err = errors.New("Ошибка при получении данных из таблицы proxy: " + err.Error())
	}
	defer rows.Close()

	if err == nil {
		proxy := new(proxyData)

		for rows.Next() {
			err = rows.Scan(&proxy.id, &proxy.ip, &proxy.port, &proxy.proxyType, &proxy.attempts)
			if err != nil {
				err = errors.New("живых прокси в БД не обнаружено " + err.Error())
				break
			}

			fmt.Println("Получен прокси", proxy.ip, proxy.port, proxy.proxyType, proxy.attempts)

			if proxy.attempts > 3 {
				_, _ = db.Exec("DELETE FROM proxies WHERE id = $1", proxy.id)
				continue
			}

			result = map[string]string{
				"ip":        proxy.ip,
				"port":      proxy.port,
				"proxyType": proxy.proxyType,
			}

			if proxyChecked(proxy.ip, proxy.port, proxy.proxyType) == true {
				fmt.Println("Прокси одобрен", proxy.ip, proxy.port, proxy.proxyType, proxy.attempts)
				break
			} else {
				fmt.Println("Прокси не дал ответ", proxy.ip, proxy.port, proxy.proxyType, proxy.attempts)
				proxy.attempts++
				r, e := db.Exec("UPDATE proxies SET attempts = $1 WHERE id = $2", proxy.attempts, proxy.id)
				fmt.Println(r, e)
				continue
			}
		}
	}

	if result == nil {
		err = errors.New("не найдено ни одного хорошего прокси в базе")
	}
	return result, err
}
