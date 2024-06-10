# luna-values-storage

## Run
```shell
    docker compose -f docker-compose.yml up -d
```

## REST API

### Получение значения
### Request

`GET /storage/values/get?id=<id>`

```
curl --location --request GET 'http://localhost:8111/storage/values/get?id=65698fabc45a676a51859302'
```

### Response

Значение найдено успешно
```http request
HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 36

{
    "id": "65698fabc45a676a51859302",
    "name": "test-value",
    "type": "string",
    "value": "string"
}
```

Значение не найдено
```http request
HTTP/1.1 404 Not Found
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 404 Not Found
    Connection: close
    Content-Type: application/json
    Content-Length: 12

{"slug": "Not found"}
```

Внутренняя ошибка
```http request
HTTP/1.1 500 Internal Server Error
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 500 Internal Server Error
    Connection: close
    Content-Type: application/json
    Content-Length: 32

{"slug": "internal-server-error"}
```

### Сохранение значения
### Request

`POST /storage/values/`

Используем multipart form-data:
#### meta-data: для метаданных значения json в формате:
```json 
        {
        "name": "new-name",
        "type": "string"  
         }
```

#### value: для самого значения (файл или строка)

```curl
curl --location 'http://localhost:8111/storage/values/set' \
--form 'value=@"/Users/danila.ivanchenko/Downloads/Combating the Sedentary Lifestyle Epidemic.pptx"' \
--form 'meta-data="{
    \"name\": \"new-name\",
    \"type\": \"string\",
    \"value\": \"new-value\"
}"'
```

### Response

Успешно сохранили
```http request
HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 36

{
    "id": "65698fabc45a676a51859302",
    "name": "new-name",
    "type": "string",
    "value": "new-value"
}
```

Внутренняя ошибка
```http request
HTTP/1.1 500 Internal Server Error
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 500 Internal Server Error
    Connection: close
    Content-Type: application/json
    Content-Length: 32

{"slug": "internal-server-error"}
```

### Удаление значения
### Request

`POST /storage/values/delete?id=65698fabc45a676a51859302`

```curl
curl --location 'http://localhost:8111/storage/values/delete?id=65698fabc45a676a51859302' \ 
```

### Response

Успешно удалили
```http request
HTTP/1.1 200 OK
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 200 OK
    Connection: close
    Content-Type: application/json
    Content-Length: 12

65698fabc45a676a51859302
```

Внутренняя ошибка
```http request
HTTP/1.1 500 Internal Server Error
    Date: Thu, 24 Feb 2011 12:36:30 GMT
    Status: 500 Internal Server Error
    Connection: close
    Content-Type: application/json
    Content-Length: 32

{"slug": "internal-server-error"}
```