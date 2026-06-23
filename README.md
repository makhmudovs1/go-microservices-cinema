# Cinema - пример микросервисов на Go с Docker, Kubernetes и MongoDB

## Обзор

Cinema — пример проекта, который показывает использование микросервисов для вымышленного кинотеатра.
Бэкенд Cinema состоит из 4 микросервисов, написанных на Go. Для хранения данных используется MongoDB, а Docker помогает изолировать и запускать всю экосистему.

 * Сервис фильмов: хранит информацию о фильмах, например название, рейтинг и другие данные.
 * Сервис сеансов: хранит информацию о расписании показов.
 * Сервис бронирований: хранит информацию о бронированиях.
 * Сервис пользователей: хранит данные пользователей и взаимодействует с другими сервисами.

Используемые контейнерные образы поддерживают несколько архитектур: `amd64`, `arm/v7` и `arm64`.

## Запуск

Приложение можно запустить на **локальной машине** через Docker Compose:

* [локальная машина (docker compose)](./docs/localhost.md)

## Как использовать сервисы Cinema

* [эндпоинты](./docs/endpoints.md)

## Архитектура

* [как устроены микросервисы, Docker Compose и Traefik](./docs/architecture.md)

## Полезные материалы

* [Микросервисы - Martin Fowler](http://martinfowler.com/articles/microservices.html)
* [Документация Traefik Proxy](https://doc.traefik.io/traefik/)
* [Go-драйвер MongoDB](https://www.mongodb.com/docs/drivers/go/current/)
* [Канал MongoDB про Go](https://www.youtube.com/c/MongoDBofficial/search?query=golang)

## Скриншоты

### Архитектура

![общая схема](docs/images/overview.jpg)

### Главная страница

![главная страница сайта](docs/images/website-home.jpg)

### Список пользователей

![страница списка пользователей](docs/images/website-users.jpg)
