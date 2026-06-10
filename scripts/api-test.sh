#!/usr/bin/env bash

BASE_URL="http://localhost:6969/api/v1"
USERNAME="testuser_$(date +%s)"
PASSWORD="testpass123"

sep() { echo; echo "=== $1 ==="; }

sep "Healthcheck"
curl -i "$BASE_URL/healthcheck"

sep "Register"
curl -i -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"$USERNAME\",\"password\":\"$PASSWORD\"}"

sep "Login + save token"
TOKEN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"name\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" \
  | jq -r '.token')
echo "TOKEN=$TOKEN"

sep "Private endpoint без токена (ожидаем 401)"
curl -i -X POST "$BASE_URL/boards" \
  -H "Content-Type: application/json" \
  -d '{"name":"JWT test board"}'

sep "Плохой токен (ожидаем 401)"
curl -i "$BASE_URL/boards/1" \
  -H "Authorization: Bearer bad-token"

sep "Create board + save id"
BOARD_ID=$(curl -s -X POST "$BASE_URL/boards" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"JWT test board"}' \
  | jq -r '.id')
echo "BOARD_ID=$BOARD_ID"

sep "Get board"
curl -i "$BASE_URL/boards/$BOARD_ID" \
  -H "Authorization: Bearer $TOKEN"

sep "Create column 1 (Todo)"
COLUMN_1_ID=$(curl -s -X POST "$BASE_URL/boards/$BOARD_ID/columns" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Todo","targetPos":0}' \
  | jq -r '.id')
echo "COLUMN_1_ID=$COLUMN_1_ID"

sep "Create column 2 (Done)"
COLUMN_2_ID=$(curl -s -X POST "$BASE_URL/boards/$BOARD_ID/columns" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Done","targetPos":1}' \
  | jq -r '.id')
echo "COLUMN_2_ID=$COLUMN_2_ID"

sep "Get columns"
curl -i "$BASE_URL/boards/$BOARD_ID/columns" \
  -H "Authorization: Bearer $TOKEN"

sep "Move column 2 на позицию 0"
curl -i -X PATCH "$BASE_URL/columns/$COLUMN_2_ID/move" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"targetPos":0}'

sep "Get columns после move"
curl -i "$BASE_URL/boards/$BOARD_ID/columns" \
  -H "Authorization: Bearer $TOKEN"

sep "Out-of-range position (ожидаем 422)"
curl -i -X POST "$BASE_URL/boards/$BOARD_ID/columns" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Invalid","targetPos":999}'

sep "Create task 1"
TASK_1_ID=$(curl -s -X POST "$BASE_URL/columns/$COLUMN_1_ID/tasks" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"First task","description":"Smoke task","targetPos":0}' \
  | jq -r '.id')
echo "TASK_1_ID=$TASK_1_ID"

sep "Create task 2"
TASK_2_ID=$(curl -s -X POST "$BASE_URL/columns/$COLUMN_1_ID/tasks" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Second task","description":"Move task","targetPos":1}' \
  | jq -r '.id')
echo "TASK_2_ID=$TASK_2_ID"

sep "Get tasks column 1"
curl -i "$BASE_URL/columns/$COLUMN_1_ID/tasks" \
  -H "Authorization: Bearer $TOKEN"

sep "Move task 2 в column 2"
curl -i -X PATCH "$BASE_URL/tasks/$TASK_2_ID/move" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"targetColumnId\":$COLUMN_2_ID,\"targetPos\":0}"

sep "Missing board (ожидаем 404)"
curl -i "$BASE_URL/boards/999999999" \
  -H "Authorization: Bearer $TOKEN"

sep "Missing column для task (ожидаем 404)"
curl -i -X POST "$BASE_URL/columns/999999999/tasks" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Missing column","description":"","targetPos":0}'

sep "Delete task 1"
curl -i -X DELETE "$BASE_URL/tasks/$TASK_1_ID" \
  -H "Authorization: Bearer $TOKEN"

sep "Delete column 1"
curl -i -X DELETE "$BASE_URL/columns/$COLUMN_1_ID" \
  -H "Authorization: Bearer $TOKEN"

sep "Delete board"
curl -i -X DELETE "$BASE_URL/boards/$BOARD_ID" \
  -H "Authorization: Bearer $TOKEN"

sep "Check deleted board (ожидаем 404)"
curl -i "$BASE_URL/boards/$BOARD_ID" \
  -H "Authorization: Bearer $TOKEN"

echo
echo "Done."