# Distributed Arithmetic Calculator

Этот проект реализует распределённый вычислитель арифметических выражений. Он состоит из двух основных компонентов:

1. **Оркестратор**: Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения.
2. **Агент**: Демон, который получает задачи от оркестратора, выполняет их и отправляет результаты обратно на сервер.

Каждая операция (сложение, вычитание, умножение, деление) имеет своё время выполнения которое задаётся в переменных среды (в файле .env). Поэтому запускать оркестратора и агента нужно ИМЕННО ТАК КАК ПОКАЗАНО В ИНСТРУКЦИИ. Также в переменных среды задаётся колличество горутин-агентов которые будут работать одновременно.

## Структура проекта

- `agent/` - Код агента.
  - `cmd/agent/` - Точка входа для запуска агента.
  - `internal/app/` - Логика работы агента.
- `orchestrator/` - Код оркестратора.
  - `cmd/calc-service/` - Точка входа для запуска оркестратора.
  - `http/server/handler/` - Обработчики HTTP-запросов.
  - `internal/app/` - Логика обработки выражений.
  - `models/` - Определения структур данных.
- `.env` - Файл для хранения переменных среды.

## Установка и запуск

### Предварительные требования

- Go 1.22
- Это не требование, но удобнее всего использовать проект будет на каком либо Linux дистрибутиве или WSL

### Установка

1. Клонируйте репозиторий:
```sh
git clone https://github.com/igicava/dac.git
```

2. Перейдите в проект:
```sh
cd dac
```

3. Установите зависимости:
```sh
go mod tidy
```

### Запуск оркестратора
1. Перейдите в каталог оркестратора:

```sh
cd orchestrator
```
2. Запустите оркестратор:

```sh
go run cmd/calc-service/main.go
```

### Запуск агента
После запуска оркестратора нужно запустить агента. Создайте 2 терминал и перейдите в КОРНЕВУЮ дирректорию проекта. Далее:
1. Перейдите в каталог агента:
```sh
cd agent
```

2. Запустите агента:
```sh
go run cmd/agent/main.go
```

## API
### Оркестратор
1. Добавление вычисления арифметического выражения
```sh
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
      "id": "<уникальный идентификатор выражения>",
      "expression": "<строка с выражением>"
}'
```

2. Получение списка выражений
```sh
curl --location 'localhost:8080/api/v1/expressions'
```

3. Получение выражения по его идентификатору
```sh
curl --location 'localhost:8080/api/v1/expressions/<id выражения>'
```

4. Получение задачи для выполнения (для агента)
```sh
curl --location 'localhost:8080/internal/task'
```

5. Прием результата обработки данных (от агента)
```sh
curl --location 'localhost:8080/internal/task' \
--header 'Content-Type: application/json' \
--data '{
      "id": "<идентификатор задачи>",
      "result": <результат>
}'
```

## Примеры использования 
1. Добавление вычисления арифметического выражения
```sh
curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
      "id": "1",
      "expression": "2 + 2 * 2"
}'
```

2. Получение списка выражений
```sh
curl --location 'localhost:8080/api/v1/expressions'
```
Вывод:
```json
{"expressions":[{"id":"1","expression":"2 + 2 * 2","status":"completed","result":6}]}
```

3. Получение выражения по его идентификатору
```sh
curl --location 'localhost:8080/api/v1/expressions/1'
```
Вывод:
```json
{"expression":{"id":"1","expression":"2 + 2 * 2","status":"completed","result":6}}
```

## Возможные ошибки:

Помните, что чаще всего возможные ошибки исправляются установкой версии Go 1.22.3, а также строгим выполнением действий в инструкции.