# Что использовал
- БД: PostgreSQL (драйвер pq)
- Роутер: chi
- Тестирование: testify и mock
- Чтение конфига: cleanenv
- Развертиывание приложения: Docker
- Тестирование API: Postman
# Структура проекта
- cmd в нем находится main файл, который инициализирует конфиг, БД, запускает сервер и обработчики
- internal вся внутренняя кухня
- internal/config конфигурационный файл и его парсер
- internal/http-server/handlers - обработчики запросов
- internal/storage - ошибки связанные с базой данных
- storage/postgres - запросы к БД
- scripts содержит shell-скрипт для ожидания поднятия БД в контейнере перед тем, как к нему подключится контейнер приложения
# Сборка и запуск проекта происходит с помощью команды   
`docker build -t go-app . && docker-compose up --build go-app`
Приложение будет готово к использованию после вывода в консоль строки "Postgres is up - executing command"
# Для запуска тестов тестов можно воспользоваться командой   
 `docker-compose run go-app go test ./.../tests`
# Что было сделано
- При инициализации БД происходит создание двух пустых таблиц(segments и clients_segments), нужных для работы API, если они еще не были созданы
- 5 методов:
   - POST-запрос для создания нового сегмена (принимает на вход название сегмента) "/seg/new"
   - Delete-запрос для удаления сегмента (принимает на вход название сегмента) "seg/del" (после удаления сегмента, сегмент также удалится у всех пользователей, и все изменения будут внесены в историю)
   - Patch-запрос для обновления информации о сегментах пользователя (на вход - добавляемые сегменты, удаляемые сегменты и ID пользователя) Я не стал делать ограничение на ID пользователя по внешнему ключу, так как не стал делать отдельную таблицу для пользователей(в задании этого написано не было + так легче выводить пользователя, если у него нет сегментов) "/user"
   - Get-запрос для получения информации о сегментах пользователя(на вход ID пользователя) "/user" 
   - Get-запрос для получения истории пользователя (на вход ID пользователя и дата в формате "YYYY-MM") на выходе ссылка на локальный файл на компьютере (Скорее всего лучшим решением для создания ссылки было бы создание файла, загрузка его на какой-то сервис, например Google Disc, и потом уже сделать ссылку от этого сервиса). Получить информацию из этого файла можно через
  `docker exec -it <ID контейнера> bash`
"/log"
- Также к этим 5 методам были написаны unit-тесты
# ПРИМЕРЫ ЗАПРОСОВ
Запросы в Postman:
1. Post-запрос на добавление сегмента
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/225e7c12-7db8-49db-879d-bb9e951940e5)
   Добавим еще сегментов : AVITO_PERFORMANCE_VAS, AVITO_DISCOUNT_30 и AVITO_DISCOUNT_50
2. Delete-запрос на удаление сегмента
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/9822fcf7-e25c-48c6-bfdf-957468fcc762)
3. Patch-запрос на обновление сегментов пользователя
  ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/d979edbe-bf6f-466b-9fe8-de1a1859a3b2)
4. Get-запрос для получения сегментов пользователя
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/b52bd9db-2d47-45aa-9467-dd6483a3ac0b)
5. Get-запрос для получения истории ссылки на файл истории пользователя
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/22ea93d7-6990-4301-8b60-487d23d99776)
   Вот что находится в этом файле:
   ![image](https://github.com/AleksNesterzz/avito_backend_task_2023/assets/109950730/b384f0eb-fa6e-49f8-8aaf-8a96d73ab153)
   Сегмент AVITO_END был добавлен и удален за кадром до создания сегмента AVITO_VOICE_MESSAGES.

   
   

   

