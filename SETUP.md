# LinkUp - Инструкции по запуску

## 🚀 Быстрый старт

### Вариант 1: Docker Compose (Рекомендуется)

1. **Клонируйте репозиторий**:
```bash
git clone <repository-url>
cd LinkUp
```

2. **Запустите все сервисы**:
```bash
docker-compose up --build
```

3. **Откройте приложение**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/index.html

### Вариант 2: Локальная разработка

#### Backend (Go)

1. **Установите зависимости**:
```bash
go mod download
```

2. **Настройте базу данных PostgreSQL**:
```bash
# Создайте базу данных
createdb linkup

# Установите переменные окружения
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=5432
export DB_NAME=chat
export JWT_SECRET=your-secret-key
```

3. **Запустите сервер**:
```bash
go run cmd/server/main.go
```

#### Frontend (React)

1. **Перейдите в папку frontend**:
```bash
cd frontend
```

2. **Установите зависимости**:
```bash
npm install
```

3. **Запустите development сервер**:
```bash
npm start
```

## 🔧 Настройка

### Переменные окружения

Создайте файл `.env` в корне проекта:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=linkup

# JWT
JWT_SECRET=your-super-secret-jwt-key

# Server
PORT=8080
HOST=localhost

# File Upload
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10485760

# 2FA
TOTP_ISSUER=LinkUp

# AI Assistant (опционально)
OPENAI_API_KEY=your-openai-key

# GitHub Integration (опционально)
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
```

### База данных

1. **Установите PostgreSQL**:
```bash
# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib

# macOS
brew install postgresql

# Windows
# Скачайте с https://www.postgresql.org/download/windows/
```

2. **Создайте базу данных**:
```bash
sudo -u postgres createdb linkup
```

3. **Создайте пользователя (опционально)**:
```bash
sudo -u postgres psql
CREATE USER linkup_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE linkup TO linkup_user;
\q
```

## 📱 Использование

### Регистрация и вход

1. Откройте http://localhost:3000
2. Нажмите "Sign up here" для регистрации
3. Заполните форму регистрации
4. Войдите в систему

### Создание комнаты

1. Нажмите кнопку "New Room" в боковой панели
2. Заполните название и slug комнаты
3. Выберите приватность комнаты
4. Нажмите "Create Room"

### Отправка сообщений

1. Выберите комнату из списка
2. Введите сообщение в поле ввода
3. Нажмите Enter или кнопку отправки

### Дополнительные функции

- **Реакции**: Наведите на сообщение и нажмите на иконку смайлика
- **Файлы**: Перетащите файл в поле ввода или нажмите на иконку скрепки
- **Поиск**: Используйте поиск для поиска сообщений
- **Настройки**: Перейдите в профиль для настройки аккаунта

## 🛠️ Разработка

### Структура проекта

```
LinkUp/
├── cmd/server/          # Точка входа приложения
├── internal/            # Внутренние пакеты
│   ├── app/            # Настройка приложения
│   ├── auth/           # Аутентификация
│   ├── handlers/       # HTTP обработчики
│   ├── models/         # Модели данных
│   ├── storage/        # Работа с БД
│   └── utils/          # Утилиты
├── frontend/           # React фронтенд
│   ├── src/
│   │   ├── components/ # React компоненты
│   │   ├── pages/      # Страницы
│   │   ├── services/  # API сервисы
│   │   └── stores/    # Состояние приложения
└── docs/              # Swagger документация
```

### Добавление новых функций

1. **Backend**:
   - Добавьте модель в `internal/models/models.go`
   - Создайте обработчики в `internal/handlers/`
   - Добавьте маршруты в `internal/app/app.go`

2. **Frontend**:
   - Создайте компоненты в `frontend/src/components/`
   - Добавьте страницы в `frontend/src/pages/`
   - Обновите типы в `frontend/src/types/index.ts`

### Тестирование

```bash
# Backend тесты
go test ./...

# Frontend тесты
cd frontend
npm test
```

## 🚀 Деплой

### Production сборка

```bash
# Frontend
cd frontend
npm run build

# Backend
go build -o linkup cmd/server/main.go
```

### Docker

```bash
# Сборка образа
docker build -t linkup .

# Запуск
docker run -p 8080:8080 linkup
```

### Docker Compose

```bash
# Production
docker-compose -f docker-compose.prod.yml up -d
```

## 🔍 Отладка

### Логи

```bash
# Docker Compose логи
docker-compose logs -f

# Конкретный сервис
docker-compose logs -f backend
```

### База данных

```bash
# Подключение к БД
psql -h localhost -U postgres -d linkup

# Просмотр таблиц
\dt

# Выход
\q
```

### WebSocket

Откройте Developer Tools в браузере и проверьте вкладку Network для WebSocket соединений.

## 📞 Поддержка

Если у вас возникли проблемы:

1. Проверьте логи приложения
2. Убедитесь, что все сервисы запущены
3. Проверьте переменные окружения
4. Создайте issue в репозитории

---

**Удачной разработки! 🚀**
