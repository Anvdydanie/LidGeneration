package model

import (
	"LidGeneration/static"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/*
Через функцию проходят все get и post запросы. В зависимости от параметров, запросы идут либо напрямую либо через прокси,
которые берутся из БД.
*/
func httpRequest(siteUrl, method string, params url.Values, contentType string, useProxy bool) (result []byte, err error) {
	var req *http.Request
	var resp *http.Response
	var timeout = time.Duration(static.REQUEST_TIMEOUT * time.Second)
	var userAgent = static.DEFAULT_USER_AGENT
	if contentType == "" {
		contentType = static.DEFAULT_CONTENT_TYPE
	}

	client := &http.Client{
		Timeout: timeout,
	}
	if useProxy == true {
		var proxyD map[string]string
		// забираем прокси из базы
		proxyD, err = getProxyFromTable()
		if err == nil {
			proxyUrl, err := url.Parse(proxyD["proxyType"] + "://" + proxyD["ip"] + ":" + proxyD["port"])
			if err == nil {
				client = &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyUrl),
					},
					Timeout: timeout,
				}
			} else {
				err = errors.New("Возникла проблема с адресом прокси: " + err.Error())
			}
		} else if err != nil && static.USE_MY_IP == true {
			// в случае отсутствия живого прокси и установленного флажка USE_MY_IP на true, используем свой ip
			err = nil
		}
	}

	if err == nil {
		req, err = http.NewRequest(strings.ToUpper(method), siteUrl, strings.NewReader(params.Encode()))
		if err == nil {
			req.Header.Set("User-Agent", userAgent)
			req.Header.Set("Content-Type", contentType)
			resp, err = client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				result, err = ioutil.ReadAll(resp.Body)
			} else {
				err = errors.New("при выполнении запроса получено " + err.Error())
			}
		}
	}

	if err != nil {
		err = errors.New("Ошибка при http запросе " + siteUrl + ". Сообщение: " + err.Error())
	}

	return result, err
}
