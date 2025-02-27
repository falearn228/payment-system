# payment-system
A store created for the send and get money transfers.

## Основные моменты

Используемые технологии -> Docker, Docker-compose, PostgreSQL, Golang (SQLC, Viper) 

`migrations` - файл с миграциями для БД \
`internal` - внутренний код проекта 

`cmd/server` - основной файл сервиса 

Почему был выбран SQLC ? 

* Типобезопасность на этапе компиляции
* Автоматическая генерация кода
* Легче поддерживать и изменять SQL-запросы
* Меньше ошибок в рантайме
* (ORM я не пользовался, познакомился на курсе по Go с SQLC и решил воспользоваться им в этом сервисе)

## Требования
- Docker
- Docker Compose

## Быстрый старт

1. Клонируйте репозиторий любым удобным спосом, пример с https:
```bash
git clone https://github.com/falearn228/payment-system.git
cd payment-system
```

В **Makefile** описаны все возможные команды. \
в **app.config.env** установлены переменные окружения для подключения к БД (также там есть адресс тестовой БД, но к сожалению, мне не хватило свободного времени, чтобы реализовать E2E и юнит тесты для использования). 

2. Собираем, скачиваем контейнеры, перейдя в папку **payment-system**
```bash
make build
```

3. Поднимаем контейнеры:
```bash
make up

# Проверяем, что все запустилось
make ps

# Остановка всех сервисов, при необходимости завершить работу
make down
```

4. Использование **API**:
```bash
Проверка баланса:
$ curl http://localhost:8080/api/wallet/<адрес_кошелька>/balance

Последние транзакции:
$ curl http://localhost:8080/api/transactions?count=5

Отправка транзакции:
$ curl -X POST -H "Content-Type: application/json" -d '{"from":"адрес_отправителя","to":"адрес_получателя","amount":"100.50"}' http://localhost:8080/api/send
```

5. Линтер для Go (Скорее всего потребуется **sudo** для установки)
```bash
sudo make lint-prepare

# Обычная проверка
make lint

# Проверка с исправлением
make lint-fix

# Быстрая проверка
make lint-fast

# Подробная проверка
make lint-verbose

# Создание JSON отчета
make lint-json
```
