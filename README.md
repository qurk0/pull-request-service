# Проект: тестовое задание Avito Tech 2025 Autumn

Потом придумаю как лучше это оформить, щас тут буду расписывать какие решения и почему я принимал.

Пока что список принятых решений:
1. Переменная LOGLEVEL в ENV-переменных, по умолчанию будет стоять LevelInfo, если вдруг захочется взять другой уровень логирования - в docker-compose.yml меняем LOGLEVEL на уровень: 'warn', 'error', 'debug' (пока что только строчные литеры работают, в дальнейшем может будет свободный регистр в переменных окружения и буду всё к lowerCase приводить, пока не до этого)

2. При создании команды если встречаем ID существующих членов команды - обновляем данные по ним: меняем username (возможно сделаю unique и не буду менять), team_name и is_active.

3. По документации есть конфликт:

Было:
  /pullRequest/reassign:
    post:
      tags: [PullRequests]
      summary: Переназначить конкретного ревьювера на другого из его команды
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [ pull_request_id, old_user_id ]
              properties:
                pull_request_id: { type: string }
                old_user_id: { type: string }
            example:
              pull_request_id: pr-1001
              old_reviewer_id: u2
У нас required-полем является old_user_id, а в example приходит old_reviewer_id. Решил оставить old_reviewer_id, считаю это логически более правильным.