# Posts Service - Ozon Test Task
Система для добавления и чтения постов и комментариев с использованием GraphQL

## Запуск проекта
Для развертывания и запуска проекта используйте:

`docker-compose up -d --build` (по умолчанию тип хранилища PostgreSQL)

`docker-compose up -d --build -e storage_type=[in-memory|postgres]` (выбор способа хранения данных)

`make unit-test-cover` (для запуска unit тестов)

API - [Postman](https://red-shadow-37838.postman.co/workspace/My-Morkspace-443094d6-ceda-482c-8ebb-1249387d97le/collection/443094d6-ceda-482c-8ebb-1249387d97le)
