# Department Service

Сервис для управления иерархией подразделений и сотрудниками.

Реализует:
- CRUD для подразделений
- вложенную структуру (дерево)
- сотрудников внутри подразделений
- перенос сотрудников при удалении департамента
- каскадное удаление
- миграции через Goose
- REST API на `net/http`

---

# ⚙️ Технологии

- Go (net/http)
- GORM
- PostgreSQL 16
- Goose (миграции)
- Docker / docker-compose

---

# 📁 Структура проекта

```
cmd/
  hitalent-test-task — точка входа в приложение
  migrator/          — сервис миграций

internal/
  config/            — загрузка конфигураций
  app/               — запуск приложения
  http/              — HTTP обработчики и request/response структуры
  service/           — бизнес-логика
  repository/        — доступ к БД
  domain/models/     — модели

migrations/          — SQL миграции goose
pkg/                 — инициализация БД
```

---

#  Запуск проекта

## 1. Клонировать репозиторий

```bash
git clone <https://github.com/Abazin97/hitalent-test-task>
cd <hitalent-test-task>
```

---

## 2. Запуск через Docker

```bash
docker-compose up --build
```

---

## 3. Сервисы

После запуска доступны:

| Сервис | URL |
|--------|-----|
| API | http://localhost:8080 |
| Postgres | localhost:5432 |

---

# Миграции

Миграции выполняются автоматически через контейнер:

```
migrator
```

Используется `goose`.

---

## Статус миграций

```
goose: no migrations to run
migrations applied successfully!
```

---

# API

### Создать департамент

```http
POST /departments
```

```json
{
  "name": "kitchen",
  "parent_id": 1
}
```

---

### Создать сотрудника

```http
POST /departments/{id}/employees
```

```json
{
  "full_name": "John Doe",
  "position": "Engineer",
  "hired_at": "2026-06-01"
}
```
---

### Получить департамент (дерево)

```http
GET /departments/{id}?depth=3&include_employees=true
```

---

### Обновить департамент

```http
PATCH /departments/{id}
```

```json
{
  "name": "New name",
  "parent_id": 2
}
```

---

### Удалить департамент

#### Cascade

```http
DELETE /departments/{id}?mode=cascade
```

#### Reassign employees

```http
DELETE /departments/{id}?mode=reassign&reassign_to_department_id=2
```


---

#  Бизнес-правила

## Department

- name: 1–200 символов
- уникален в пределах одного parent
- нельзя создать цикл в дереве
- parent_id должен существовать

---

## Employee

- full_name: 1–200
- position: 1–200
- создаётся только в существующем департаменте

---

## Удаление департамента

### cascade
- удаляет департамент
- всех сотрудников
- всех потомков

### reassign
- сотрудники переводятся в другой департамент
- затем департамент удаляется

---

# 🐳 Docker Compose

```yaml
services:
  db:
    image: postgres:16

  migrator:
    build: .
    command: ./migrator

  app:
    build: .
    command: ./app
```

---
