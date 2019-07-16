# LidGeneration
The service was created to simplify the analysis of competition in the sphere of services and sales on the Internet in Russian cities. Using this service user gets relevant search results from 5 services (yandex, google, vk.com, avito, yandex advert) via 1 button click. It's usefull for people who wants to analyse online market competition.

Project uses 2 third-party packages:

github.com/PuerkitoBio/goquery

github.com/mattn/go-sqlite3


Algorithm performance criteria:
1. You type theme of service.
2. You choose city in Russia.
3. Choose type of search. Sample of results can be strict or soft. In first case in title or description of parsed results there must match all keywords from theme field. In second case there must be at least 1 keyword match with title or description.
4. Services to parse are: yandex and google search, yandex advertisment, vk.com groups, avito offers. The numbers of results are limited to top 30.
5. When parsing is done, special algorithm analyse list of results on keywords and their synonyms (using yandex vocabulary service). If there is match, depending on chose type of search than result is considered relevant and will be shown to user.
6. When backend work is done user wil see 5 tables with results (title, link, description).
6. User can choose to show him all results or only relevant.
7. User can although delete any row in table he wants.
