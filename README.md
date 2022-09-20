# Balance serice

  

**Микросервис для работы с балансом пользователей**

Для запуска сервиса в докере введите команду `make docker-start`
Для запуска локально использовать команду `make run`
Файл для конфигурации распологается в `./k8s/values_local.yaml`
Ручное тестирование сервиса:
* Для тестирования можно использовать *swaggerUI*, располагается по адресу `http://localhost:8000/swagger/index.html`
* Для тестирования через curl или postman, базовый url api располагается по адресу
`http://localhost:8000/`
* Для запуска модульных и интеграционных тестов используйте команду 
`make test`
в тестах используется подключение к бд postgres замените строку поключения с ключем **pg_dsn** в файле конфигурации   

## Маршруты
 

### GET /balance?user_id="user_1"&currency="EUR"

Получить баланс пользователя в валюте.
*user_id* - обязательный праметр
*currency* - не обязательный(по умолчанию "RUB") формат ISO 4217 ("RUB", "EUR", "USD")

### POST /transfer

Body "{"*amount*":1000, "*sender_id*":"1", "*receiver_id*":"2"}"
*amount* - перечисляемые средства в копейках

Перевести средства с одного счета на другой
* Пополнение баланса
Для пополнения баланса sender_id оставить пустым или использовать нулевой *uuid*
* Снятие со счета
Для снятия receiver_id оставить пустым или использовать нулевой uuid
* Перевод средств
Для перевода поля receiver_id и sender_id должны быть заполнены.

### GET /transfers
    http://localhost:8000/transfers?amountEQ=1&amountGTE=2&amountLTE=3&dateFrom=4&dateTo=5&page=6&perPage=7&transferType=8&userID=9

Параметры маршрута:
1. amountEQ int64
2. amountGTE int64
3. amountLTE int64
4. dateFrom string (RFC3999)
5. dateTo string (RFC3999)
6. page int
7. perPage int
8. transferType string("debit", "credit", "transfer")
9. userID string 