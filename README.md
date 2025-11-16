# Проект: тестовое задание Avito Tech 2025 Autumn

Данный проект является реализацией сервиса для управления PullRequest'ами.

## Стек технологий, которые использованы в данном проекте:
1. **Язык программирования:** Go
- - Библиотеки: cleanenv, fiber, jackx/pgx/v5
2. **База данных:** PostgreSQL (+ миграции с помощью Goose)
3. **Контейнеризация с помощью Docker:** на момент создания коммита нет ничего, кроме директории Docker. Потом напишу 3 файла: докер компоуз для развертывания всего проекта и 2 докерфайла для сборки сервиса и мигратора

## Эндпоинты, которые я не знаю - стоит заносить в документацию или нет, но сказать о них надо
1. /metrics - возвращает метрики. Тут я старался использовать prometheus-метрики с кастомными счётчиками в пакете internal/metrics

## Принятые решения в рамках разработки проекта
1. Переменная LOGLEVEL в ENV-переменных, по умолчанию будет стоять LevelInfo, если вдруг захочется взять другой уровень логирования - в docker-compose.yml меняем LOGLEVEL на уровень: 'warn', 'error', 'debug' (пока что только строчные литеры работают, в дальнейшем может будет свободный регистр в переменных окружения и буду всё к lowerCase приводить, пока не до этого)

2. При создании команды если встречаем ID существующих членов команды - обновляем данные по ним: меняем username (возможно сделаю unique и не буду менять), team_name и is_active.

### Нагрузочное тестирование

Проводил с помощью hey:

hey -n 20 -c 5 \
  -m POST \
  -H "Content-Type: application/json" \
  -d '{
    "team_name": "backend",
    "members": [
      {
        "user_id": "u1",
        "username": "Alice",
        "is_active": true
      },
      {
        "user_id": "u2",
        "username": "Bob",
        "is_active": true
      },
      {
        "user_id": "u3",
        "username": "Charlie",
        "is_active": true
      }
    ]
  }' \
  http://localhost:8080/team/add

hey -n 20 -c 5 \
  "http://localhost:8080/team/get?team_name=backend"

hey -n 20 -c 5 \
  -m POST \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "u2",
    "is_active": false
  }' \
  http://localhost:8080/users/setIsActive          

### E2E тестирование

В директории app лежит файл для e2e тестирования. Мои результаты:


=== RUN   TestPRLifecycleE2E
    e2e_test.go:42: step 1: /team/add
    e2e_test.go:67: step 2: /pullRequest/create
    e2e_test.go:96: step3: /users/getReview
    e2e_test.go:117: step 4: /pullRequest/reassign
    e2e_test.go:144: step 5: /pullRequest/merge
    e2e_test.go:160: step 6: /pullRequest/reassign on merged PR (expect 409 PR_MERGED)
--- PASS: TestPRLifecycleE2E (0.02s)
PASS
ok      github.com/qurk0/pr-service/e2e 0.019s
