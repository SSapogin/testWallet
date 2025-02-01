- Переименовываем .env.example в .env
- Запускаем `docker compose up -d`
- Тесты `go test -v ./...`
- с учетом покрытия `go test ./... -cover`
- Проводим нагрузочное тестирование:
`hey -z 10s -c 200 -q 1000 -m POST \
  -T "application/json" \
  -d '{"walletId":"b26d64b8-7e31-49bf-a940-bb91cc0969b4", "operationType":"DEPOSIT", "amount":1}' \
  http://localhost:8080/api/v1/wallet
`