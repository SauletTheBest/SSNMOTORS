# E-Commerce Microservices

Этот проект состоит из трёх микросервисов:

* **User Service** — управление пользователями
* **Inventory Service** — управление товарами
* **Order Service** — оформление и список заказов

Все микросервисы используют **Go**, **gRPC**, **MongoDB** и **Redis**. Реализованы по принципам **чистой архитектуры** (Clean Architecture).

---

## 🧑‍💼 User Service

### Функции:

* Регистрация и аутентификация
* Получение и обновление профиля
* Кэширование профиля пользователя по ID (Redis)

### Кэш:

* Ключ: `user_profile:<userID>`
* TTL: 5 минут

### Проверка через Redis CLI:

```bash
GET user_profile:<userID>
TTL user_profile:<userID>
```

---

## 📦 Inventory Service

### Функции:

* Получение информации о товаре по ID
* Добавление, обновление и удаление товаров
* Кэширование товара по ID (Redis)

### Кэш:

* Ключ: `product:<productID>`
* TTL: 5 минут

### Проверка через Redis CLI:

```bash
GET product:<productID>
TTL product:<productID>
```

---

## 🛒 Order Service

### Функции:

* Получение заказов по ID пользователя
* Кэширование списка заказов пользователя (Redis)

### Кэш:

* Ключ: `user_orders:<userID>`
* TTL: 5 минут

### Проверка через Redis CLI:

```bash
GET user_orders:<userID>
TTL user_orders:<userID>
```

---

## 🔗 Общие технологии

* **gRPC** — коммуникация между микросервисами
* **MongoDB** — хранение данных
* **Redis** — кэширование
* **Go** — основной язык разработки

### Запуск Redis и Mongo:

```bash
docker run -d -p 6379:6379 redis
docker run -d -p 27017:27017 mongo
```

---

## 📦 Структура микросервиса

```
├── cmd/            # main.go
├── internal/
│   ├── handler/    # gRPC endpoints
│   ├── usecase/    # бизнес-логика
│   ├── repository/ # mongo
│   ├── model/
│   └── cache/      # Redis client
├── proto/
├── config/
└── README.md
```
