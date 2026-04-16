# hackathon2026-rebuild
Переосмысление уже написанного кода для кейса на хакатоне 2026

# Запуск
1. Задать переменные окружения в `.env` (пример в `.env.example`)
2. Написать `infra/garage.toml` (пример в `infra/garage-example.toml`)
3. `docker compose up --build`

# Тестирование
В `backend/test/fakes.go` распологаются все фейковые реализации интерфейсов.

## Юнит тесты
В `backend/test/usecases_test.go` приведены юнит-тесты.
Для запуска перейти в `backend`: `go test -v -short ./test`

## End-to-end
Для запуска перейти в `backend`: `go test -tags=e2e ./test`

### Тестовое окружение
Для e2e теста требуется тестовое окружение. В `docker-compose.yaml` есть пример тестового окружения. Запуск: `docker compose --profile test up test-db test-garage`
