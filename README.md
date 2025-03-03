# Простой веб-сервер на Golang

Зависимости: Golang 1.23.6 ИЛИ Docker

Стек: Go, Echo

## Архитектура:

Handler: HTTP

Service:

Reposytory: data.txt

## Endpoints:

- POST http://host:port/api?line= - принимает на вход строку, записывает ее в файл data.txt

- GET http://host:port/api/ - инвертирует имеющиеся в data.txt строки и возвращает

## Запуск:

Manual:

```sh
go mod download

go run cmd/app/main.go

# OR

make run
```

Docker (compose):

```
docker compose -f ./compose.yaml up

# OR

make compose-run
```

Shutdown:

```sh
^C

# OR

docker compose -f ./compose.yaml down

# OR for docker compose

make compose-down

```
