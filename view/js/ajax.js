document.addEventListener("DOMContentLoaded", ready);

let searchResults = [];

function ready() {
    // поиск
    const inputTheme = document.getElementById("inputTheme");
    const inputCity = document.getElementById("inputCity");
    const citiesList = document.getElementById("citiesList");
    const isStrictSearch = document.getElementById("isStrictSearch");
    const button = document.getElementById("buttonAjaxRequest");
    // таблицы
    const yandexSearch = document.getElementById("yandex-search").getElementsByTagName('tbody')[0];
    const yandexAdvert = document.getElementById("yandex-advert").getElementsByTagName('tbody')[0];
    const googleSearch = document.getElementById("google-search").getElementsByTagName('tbody')[0];
    const vkGroups = document.getElementById("vk-groups").getElementsByTagName('tbody')[0];
    const avitoList = document.getElementById("avito-offers").getElementsByTagName('tbody')[0];
    // массивы и переменные
    let lastCityOptions = [];
    let warningMsg = document.getElementById("warning");

    // Проверяем тему
    inputTheme.onkeyup = function() {
        validate(inputTheme);
    };

    // получаем список городов
    inputCity.onkeyup = function() {
        // проверяем ввод данных
        validate(inputCity);
        // Получаем автокомплит
        if (inputCity.value.length > 1) {
            let result = sendAjaxRequest(inputCity.value, "/getCitiesList");
            if (result.status === 200 && result.response !== null ) {
                citiesList.innerHTML = "";
                result.response.forEach(function (cityName) {
                    let option = document.createElement('option');
                    option.value = cityName;
                    citiesList.appendChild(option);
                });
                // Сохраняем последний вариант автокомплита для проверки перед парсингом
                lastCityOptions = result.response;
                // Если город определен очищаем автокомплит
                if ( lastCityOptions.includes(inputCity.value) ) {
                    citiesList.innerHTML = "";
                    warningMsg.innerHTML = "";
                }
            }
        }
    };

    // отправляем запрос на парсинг выдачи
    button.onclick = function(event) {
        event.preventDefault();
        warningMsg.innerHTML = "";
        // валидиция темы, перед отправкой запроса
        if (inputTheme.value.trim().length < 2) {
            warningMsg.innerHTML = "Вам необходимо указать тему!";
        }
        // валидиция введенного города, перед отправкой запроса
        else if (
            lastCityOptions.length === 0
            ||
            !lastCityOptions.includes(inputCity.value)
        ) {
            warningMsg.innerHTML = "Не знаю такого города!";
        }
        // если все впорядке, отправляем запрос
        else {
            // Получаем данные формы
            let theme = inputTheme.value;
            let cityName = inputCity.value;
            let strictSearch = isStrictSearch.checked;
            let reqString = JSON.stringify({
                "theme": theme,
                "cityName": cityName,
                "strictSearch": strictSearch
            });
            // Отправляем запрос
            let result = sendAjaxRequest(reqString, "/parseSearchEngines");
            if (result.status === 200) {
                searchResults = [
                    {
                        response: [result.response.yandexSearch],
                        selector: yandexSearch
                    },
                    {
                        response: [result.response.yandexAdvert],
                        selector: yandexAdvert
                    },
                    {
                        response: [result.response.googleSearch],
                        selector: googleSearch
                    },
                    {
                        response: [result.response.vkGroups],
                        selector: vkGroups
                    },
                    {
                        response: [result.response.avito],
                        selector: avitoList
                    },
                ];
                showResultInTables(searchResults, "filteredResponse");
            } else {
                warningMsg.innerHTML = "К сожалению на сервере произошла ошибка: " + result.status;
            }
        }
    };
}


function validate(inputText) {
    let validRange = /^[А-Яа-я]+|(-+)$/;
    let isValid = inputText.value.match(validRange);
    if (!isValid) {
        inputText.value = "";
        inputText.placeholder = "Вводить только кириллицу";
    }
}


function sendAjaxRequest(paramsString, method) {
    const xhr = new XMLHttpRequest();
    xhr.open("POST", method, false);
    xhr.send(paramsString);
    if (xhr.status === 200) {
        try {
            let result = JSON.parse(xhr.response);
            return {status: xhr.status, response: result};
        } catch (e) {
            return {status: xhr.status, response: xhr.response}
        }
    } else {
        return {status: xhr.status, response: xhr.response}
    }
}


function showResultInTables(seArr, resultToShow) {
    seArr.forEach(function(searchEngine) {
        if (
            typeof searchEngine.response !== "undefined"
            &&
            searchEngine.response !== null
            &&
            searchEngine.response[0] !== null
            &&
            typeof searchEngine.response[0] !== "undefined"
            &&
            Array.isArray(searchEngine.response[0][resultToShow])
            &&
            searchEngine.response[0][resultToShow].length > 0
        ) {
            searchEngine.selector.innerHTML = "";
            // показываем результаты
            searchEngine.response[0][resultToShow].forEach(function (item) {
                let title = "<span class='title'>"+item.title+"</span>";
                let url = "<p class='url'><a href=\""+item.fullUrl+"\" target='_blank'>"+item.domain+"</a></p>";
                let description = item.description !== undefined ? "<p class='description'>"+item.description+"</p>" : "";
                let position = item.position !== undefined ? "<td class='position'>"+item.position+"</td>" : "";
                let deletePos = "<td id='delete-pos' onclick='deletePosition(this.parentNode)'><span> X </span></td>";
                searchEngine.selector.innerHTML += "<tr>" + position + "<td>" + title + url + description + "</td>" + deletePos + "</tr>";
            })
        } else {
            // TODO сообщаем причину ошибки
            searchEngine.selector.innerHTML += "<tr><td>Не найдено результатов, отвечающих требованиям поиска.</td></tr>"
        }
    });
}

// показывает полный результат
function showFullResults() {
    showResultInTables(searchResults, "fullResponse")
}

// показывает фильтрованный результат
function showFilteredResults() {
    showResultInTables(searchResults, "filteredResponse")
}

// удаляет строку
function deletePosition(elem) {
    elem.style.display = "none";
}