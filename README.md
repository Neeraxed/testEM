# testEM

## Локальная сборка
Собираем бинарник:
```bash
make build
```

Запуск:
```bash
make start
```

Генерация swagger:
```bash
make docs
```

## Запуск в Docker
Сборка образа
```bash
docker build -t test_em .
```
Запуск контейнера
```bash
docker container run -p 3333:3333 -ePOSTGRESDSN="host=localhost user=tester password=pswd dbname=test sslmode=disable" -eEXTERNALURL="url" -ePORT=":3333" --net=host -t --rm test_em
```