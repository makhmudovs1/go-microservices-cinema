# Архитектура Cinema

Этот документ объясняет, из чего состоит приложение, как микросервисы общаются между собой и зачем здесь нужен Traefik.

## Общая Схема

```text
Браузер
  |
  | http://localhost
  v
Traefik container
  |
  | PathPrefix(`/`)
  v
website container
  |
  | HTTP-запросы внутри Docker-сети
  | http://users:4000/api/users/
  | http://movies:4000/api/movies/
  | http://showtimes:4000/api/showtimes/
  | http://bookings:4000/api/bookings/
  v
users / movies / showtimes / bookings containers
  |
  | MongoDB connection
  v
db container
```

## Что Такое Микросервис

Микросервис в этом проекте - это отдельное маленькое приложение, которое отвечает за одну часть системы.

В проекте есть такие сервисы:

```text
website    - UI, HTML-страницы
users      - пользователи
movies     - фильмы
showtimes  - сеансы
bookings   - бронирования
db         - MongoDB
proxy      - Traefik
```

Например, `movies` - отдельный Go-сервис. У него есть свой код, свой `Dockerfile`, свой `go.mod`, свои handlers и своя модель данных. Он отвечает на HTTP API:

```text
GET    /api/movies/
GET    /api/movies/{id}
POST   /api/movies/
DELETE /api/movies/{id}
```

Физически каждый микросервис запускается как отдельный Docker-контейнер.

## Что Делает Docker Compose

Когда запускается команда:

```bash
docker compose up --build --detach
```

Docker Compose читает `compose.yaml`, собирает образы из `Dockerfile` и запускает контейнеры.

Например:

```yaml
movies:
  build: ./movies
  image: cinema-movies:local
  command:
    - "-mongoURI"
    - "mongodb://db:27017/"
```

Это значит:

```text
1. зайти в папку ./movies
2. собрать Docker image по movies/Dockerfile
3. назвать image cinema-movies:local
4. запустить контейнер movies
5. передать ему mongoURI mongodb://db:27017/
```

`db` - это имя MongoDB-сервиса в Docker Compose. Внутри Docker-сети контейнеры могут обращаться друг к другу по имени сервиса:

```text
movies  -> mongodb://db:27017
website -> http://movies:4000/api/movies/
```

## Как Общаются Сервисы

Сервисы общаются по HTTP внутри общей Docker-сети.

Например, пользователь открывает список фильмов:

```text
http://localhost/movies/list
```

Путь запроса:

```text
Браузер
  |
  v
Traefik
  |
  v
website
  |
  | GET http://movies:4000/api/movies/
  v
movies
  |
  | запрос в MongoDB
  v
db
```

`movies` возвращает JSON с фильмами. `website` берет этот JSON, подставляет данные в HTML-шаблон и возвращает готовую страницу браузеру.

## Роль Website

`website` - это тоже отдельный сервис, но он не хранит данные. Он работает как frontend/backend-for-frontend:

```text
website:
  - принимает запросы от браузера
  - рисует HTML
  - дергает API других сервисов
  - собирает данные в страницы
```

Например, страница бронирований может требовать данные сразу из нескольких сервисов:

```text
website -> bookings
website -> users
website -> showtimes
website -> movies
```

Потому что чтобы показать бронь, нужны пользователь, сеанс и выбранные фильмы.

## Что Делает Traefik

Traefik - это reverse proxy, примерно как Nginx.

Он является первой точкой входа:

```text
Браузер -> Traefik -> нужный контейнер
```

Traefik слушает внешний порт `80`:

```yaml
ports:
  - target: 80
    published: 80
```

Поэтому сайт открывается просто через:

```text
http://localhost
```

Это то же самое, что:

```text
http://localhost:80
```

Dashboard Traefik открыт отдельно:

```text
http://localhost:8080/dashboard/#/
```

## Как Traefik Понимает, Куда Отправлять Запросы

Traefik подключен через `compose.yaml`.

Эти команды говорят Traefik читать Docker-контейнеры и слушать HTTP-вход:

```yaml
- "--providers.docker"
- "--entrypoints.web.address=:80"
```

А эти labels у `website` говорят:

```yaml
labels:
  - "traefik.http.routers.website.rule=PathPrefix(`/`)"
  - "traefik.http.services.website.loadbalancer.server.port=8000"
```

Смысл:

```text
Все HTTP-запросы, начинающиеся с `/`,
отправлять в контейнер website на порт 8000.
```

Например:

```text
http://localhost/movies/list
```

идет так:

```text
Браузер -> Traefik :80 -> website :8000
```

Потом уже `website` сам делает внутренний HTTP-запрос:

```text
website -> http://movies:4000/api/movies/
```

## Почему Сейчас Traefik Почти Лишний

В текущей конфигурации Traefik в основном отправляет все запросы в `website`:

```yaml
- "traefik.http.routers.website.rule=PathPrefix(`/`)"
```

То есть схема сейчас такая:

```text
Браузер
  |
  v
Traefik
  |
  v
website
  |
  v
users / movies / showtimes / bookings
```

В таком режиме можно было бы обойтись без Traefik:

```text
Браузер -> website:8000 -> API-сервисы -> MongoDB
```

И это как раз работало через:

```text
http://localhost:8000
```

## Когда Traefik Становится Нужен

Traefik становится полезным, когда через него идут разные маршруты в разные сервисы:

```text
Браузер / API-клиент
  |
  v
Traefik
  |
  +-- /               -> website
  +-- /api/users/     -> users
  +-- /api/movies/    -> movies
  +-- /api/showtimes/ -> showtimes
  +-- /api/bookings/  -> bookings
```

Тогда запрос:

```text
http://localhost/api/movies/
```

пойдет напрямую в `movies`, а запрос:

```text
http://localhost/users/list
```

пойдет в `website`.

Сейчас labels для API-сервисов в `compose.yaml` закомментированы:

```yaml
# - "traefik.http.routers.movies.rule=PathPrefix(`/api/movies/`)"
# - "traefik.http.services.movies.loadbalancer.server.port=4000"
```

Если раскомментировать такие labels для `movies`, `users`, `showtimes` и `bookings`, Traefik станет полноценным API Gateway.

## Важно Про Website

`website` не перенаправляет браузер на микросервисы. Он сам, на серверной стороне, делает HTTP-запросы к ним.

Браузер этого не видит. Браузер получает уже готовую HTML-страницу.

```text
Браузер -> website -> микросервисы -> website -> HTML -> Браузер
```

## Итог

Текущая схема:

```text
Браузер -> Traefik -> website -> API-сервисы -> MongoDB
```

Упрощенная схема без Traefik:

```text
Браузер -> website -> API-сервисы -> MongoDB
```

Более правильная gateway-схема:

```text
Браузер / API-клиент -> Traefik -> website или нужный API-сервис
```

Traefik в этом проекте нужен как демонстрация подхода, который часто используют в production:

```text
- единая точка входа
- маршрутизация запросов
- dashboard
- middleware
- балансировка
- TLS/HTTPS в реальных проектах
- автоматическое обнаружение сервисов через Docker labels
```
