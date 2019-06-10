package static

// используются для сравнения с переменной typeOfSearch. Указывает, что необходимо парсить: поисковую выдачу или рекламную
const SEARCH_TYPE = "Search"
const ADVERT_TYPE = "Advertising"

// используются при получении списка городов для автокомплита
const CITIES_NAME_API_URL = "http://kladr-api.ru/api.php?"
const CITIES_NAME_API_SETTLE_TYPE = "city" // city, country, ...
const CITIES_NAME_API_TYPE_CODE = "1"      // 1 = город, 2 = поселок, 4 = деревня
const SEARCH_COUNTRY_NAME = "Россия"

// api для получения id городов в таблице яндекса
const YANDEX_CITIES_ID_URL = "https://pogoda.yandex.ru/static/cities.xml"

// url для поисковой выдачи
const YANDEX_SEARCH_URL = "https://yandex.ru/search/?"

// url для рекламной выдачи
const YANDEX_ADVERT_URL = "https://yandex.ru/search/ads?"

// api и данные для получения поисковой выдачи яндекса через xml сервис
const YANDEX_SEARCH_API_XML_URL = "https://yandex.ru/search/xml?"
const YANDEX_SEARCH_API_XML_USER = "andydanie"
const YANDEX_SEARCH_API_XML_KEY = "03.69574483:6ed60187c4c45f07b37224c6a29bf84c"
const YANDEX_SEARCH_API_XML_SORT = "rlv"                                                 // только релевантная выдача
const YANDEX_SEARCH_API_XML_FILTER = "strict"                                            // строгий фильтр, без порно и прочей важной для жизни фигни
const YANDEX_SEARCH_API_XML_GROUP = "attr=d.mode=deep.groups-on-page=30.docs-in-group=1" // 50 результатов на страницу
// api и данные для получения поисковой выдачи гугла через json сервис
const GOOGLE_SEARCH_API_JSON_URL = "https://www.googleapis.com/customsearch/v1?"
const GOOGLE_API_JSON_KEY = "AIzaSyA0Y3X40Qzp0ySTaNgSmEjsIo_KJhLlPUY"
const GOOGLE_API_JSON_SID = "010786294839706012062:l4bjixkzqua"

// стандартные данные для отправки http запросов
const DEFAULT_CONTENT_TYPE = "application/x-www-form-urlencoded"
const DEFAULT_USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36"

// api и данные для получения синонимов слов
const DICTIONARY_API = "https://dictionary.yandex.net/api/v1/dicservice.json/lookup?"
const DICTIONARY_KEY = "dict.1.1.20190531T131207Z.f9f5b015b61579b9.9273c50622e75f6daf65993061321a0cb270a8f8"
const DICTIONARY_LANG = "ru-ru"

// путь к файлам и папкам
const VIEW_FOLDER_PATH = "./LidGeneration/view"
const DATABASE_FILE_PATH = "./LidGeneration/database/lidgen.db"

// количество страниц, которое необходимо спарсить в выдаче
const MAX_PAGES_TO_PARSE = 3

// таймаут на ожидание ответа при http запросе
const REQUEST_TIMEOUT = 20

// разрешаем или запрещаем использовать свой ip для парсинга, в случае если нет живых прокси
const USE_MY_IP = true

// api и данные для получения групп Вконтакте
const VK_API_GET_GROUPS = "https://api.vk.com/method/groups.search?"
const VK_API_GROUP_TYPE = "group" // Возможные значения: group, page, event
const VK_API_GROUPS_COUNT = "50"  // сколько групп возвращать в результатах
const VK_API_TOKEN = "9ab47300dd5ce2a9a73bc6bed0324fdf9eb7714eb933484e9ce2f806bf72257094391c142cb8887cd0331"
const VK_API_VERSION = "5.95"
const VK_API_GET_CITIES_ID = "https://vk.com/select_ajax.php?act=a_get_cities&country=1&str="
const VK_GROUP_URL = "https://vk.com/public"

// url и данные для парсинга авито
const AVITO_SEARCH_URL = "https://www.avito.ru/"
const AVITO_CATEGORY_ID = "3" // это параметр категории для поиска. По умолчанию 3, то есть ищем во всех категориях
