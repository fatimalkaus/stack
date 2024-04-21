# Stack

## Запуск
1. `make build`
2. `make docker-build`
3. `docker-compose up -d`
4. `make migrate-up`

## API

### Добавить элемент в стек
- Endpoint: /
- Method: POST
- Request Body: 
  {
    "data": "new item"
  }
- Response: 
  - Status Code: 204 No Content
  - Body: None 


### Посмотреть верхний элемент стека
- Endpoint: /
- Method: GET
- Request Body: None
- Response: 
  - Status Code: 200 OK
  - Body: 
    {
      "data": "top item"
    }

### Достать верхний элемент стека
- Endpoint: /
- Method: DELETE
- Request Body: None
- Response: 
  - Status Code: 200 OK
  - Body: 
    {
      "data": "popped item"
    }
