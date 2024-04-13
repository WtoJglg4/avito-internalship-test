## Запуск

1. Запустите скрипт:

```sh
./run.sh
```

## Ход решения

Как организовать хранение?
1. Один баннер - одна фича, одна фича - несколько баннеров = связь один ко многим
2. Один баннер - несколько тегов, один тег - несколько баннеров = связь многие ко многим
3. Для связи один ко многим назначим feature_id в таблице banners, как foreign key, который будет ссылаться на столбец id в features
4. Для связи многие ко многим создадим отдельную таблицу banner_tags
5. Индексация и оптимизация: Для ускорения поиска и избежания дублирования информации рекомендуется использовать индексы, особенно на полях feature_id в таблице banners и (banner_id, tag_id) в таблице banner_tags.

Кэширование
Будем использовать кеширование для разгрузки базы:
1. User запрашивает баннер
2. Если нет флага use_last_revision, идем в Redis, проверяем наличие баннера по сочетанию фичи и тега, а так же TTL.
3. Если есть флаг use_last_revision, или данных нет в Redis, или TTL > 5 минут, идем сразу в бд.
    3.1 Добавляем полученные из БД данные в Redis.
4. Возвращаем данные user`y.

Версионность баннеров(доп требование)?
Будем хранить 3 предыдущие версии баннера в отдельной таблице banner_versions, которая будет хранить номер версии, дату изменения и сам баннер.
Изменим таблицу banners - добавим новый стобец version, который будет отражать текущую версию баннера.
Логика версионирования: 
1. При каждом изменении баннера сохраняем текущую версию баннера в таблицу banner_versions
2. Удаляем самую старую версию
3. Обновляем текующую версию баннера, инкрементировав счетчик версий.
4. Закешируем новую версию баннера, удалив старый кеш, если он был.
Изменим API так, чтобы можно было чтобы можно было:
1. Просмотреть существующие версии баннера: GET /banner/{id}/versions
2. Выбрать подходящую версию: GET /banner/{id}/versions/{version}

Метод удаления баннеров по фиче или тегу, время ответа <= 100мс? Механизм отложенных действий?
Данная задача может быть потенциально непростой из-за времени ответа. Нагрузка на базу может кратно возрасти из-за операций удаления и мы не получим 100мс. Воспользуемся системой управления очередями задач, а именно будем использовать RabbitMQ.
Он позволит нам распараллелить нагрузку по воркерам + сможем получать подтверждение выполнения операции.
Логика удаления:
1. Запрос на эндпоит удаления с id фичи или тега.
2. Формируем сообщение для RabbitMQ из полученных данных.
3. Кладем сообщение в очередь
4. Воркер подхватывает сообщение из очереди и делает транзакции на удаление
5. Лог
Новый эндпоинт для удаления: DELETE: /banner

Линтер
[golintci](https://github.com/golangci/golangci-lint)

Регистрация и авторизация:
В API добавлено 2 эндпоинта:
Регистрация: /auth/sign-up
Авторизация: /auth/sign-in

Распределение ролей: admin или default
Так как нам необходимо распределить роли пользователей, для удобства решил сделать так:
Админ - один пользователь с логином и паролем, описанным в файле .env. Соответственно, админский токен можно получить, введя только этот пароль и логин. Остальным пользователям по умолчанию присваивается роль default. 

Как реализовать добавление баннера?
Этапы добавления нового баннера:
1. Убедимся, что переданные id фичи и тегов все существуют в таблицах features и tags, соответственно
2. Найдем прошлую версию баннера(если таковая есть), обновим поля is_active и updated_at(найдем по уникальному сочетанию фичи и тегов)
Как найти прошлую версию?
Решено добавить столбец с хешом тегов, т.к. это значительно ускорит поиск предыдущей версии баннера.

