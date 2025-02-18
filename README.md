Логика реализовано полностью. Все остальное in progress...

### TODO

- Тесты
- Makefile
- Docker-compose
- Нагрузочное тестирование
- Линтер
- Расписанный README

### Проблемы и их решение

1. Где реализовывать логику auth? Сервисный слой, транспортный слой --> **Решил** создать отдельный сервис отвечающий за это. 
В его ответвественности создавать токены, парсить токены, создавать если не существует пользователей в db.
2. Как иметь транзакции и не добавлять бизнес-логику в слой репозитория ? --> **Решил** -- Менеджер транзакций в GoLang / Илья Сергунин (Авито) (https://www.youtube.com/watch?v=fcdckM5sUxA&ab_channel=HighLoadChannel)
3. Что использовать как primary key? --> пришел к тому что лучше использовать username несмотря на то, что если бы мы пользовались id - bigserial, uuid индексы в postgres работали быстрее.
Потому что постоянный маппинг например для запросов инфо (coin history) сьедал бы слишком много ресурсов.