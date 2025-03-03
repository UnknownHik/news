# Новости

Сервис позволяет обновлять новости (заголовки, содержимое и категории) и 
получать пагинированный список всех новостей с возможностью указания размера страницы и 
номера страницы.

## Технологии

- **Язык:** Go (1.23)
- **Фреймворки:** 
  - Fiber: используется для создания HTTP-сервиса
  - Reform: ORM для работы с PostgreSQL
- **База данных:** PostgreSQL
- **Аутентификация:** JWT
- **Контейнеризация:** Docker Compose

## Запуск проекта

### 1. Клонирование репозитория

```bash
git clone https://github.com/UnknownHik/Merch-store
```

### 2. Запуск через Docker Compose

```bash
docker-compose up --build
```
### 3. Cервис будет доступен по адресу http://localhost:8080