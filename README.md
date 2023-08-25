# Простой веб-сервер
Небольшой веб-сервер с возможностью принимать GET и POST запросы по пути `localhost:8080/person/`
## Используемые технологии
     http, log, json
## Пример запросов 
```shell
curl -iX GET -s localhost:8080/person/
```
```shell
curl -iX POST -s localhost:8080/person/ -H 'Content-Type: application/json' -d '{"name":"Leo","age":26}'
```
