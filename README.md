# Инкремент 10
Сервер:
  Добавьте функциональность подключения к базе данных. В качестве СУБД используйте PostgreSQL не ниже 10 версии.
  Добавьте в сервер хендлер GET /ping, который при запросе проверяет соединение с базой данных. При успешной проверке хендлер 
  должен вернуть HTTP-статус 200 OK, при неуспешной — 500 Internal Server Error.
  Строка с адресом подключения к БД должна получаться из переменной окружения DATABASE_DSN или флага командной строки -d.
  Для работы с БД используйте один из следующих пакетов:
  - database/sql,
  - github.com/jackc/pgx,
  - github.com/lib/pq,
  - github.com/jmoiron/sqlx.

# Инкремент 11
Перепишите сервер для сбора метрик таким образом, чтобы СУБД PostgreSQL стала хранилищем метрик вместо текущей реализации.
Сервису нужно самостоятельно создать все необходимые таблицы в базе данных. Схема и формат хранения остаются на ваше усмотрение.
Для хранения значений gauge рекомендуется использовать тип double precision.
При отсутствии переменной окружения DATABASE_DSN или флага командной строки -d или при их пустых значениях вернитесь последовательно к:
хранению метрик в файле при наличии соответствующей переменной окружения или флага командной строки;
хранению метрик в памяти.

# Инкремент 12
Сервер:
  Добавьте новый хендлер POST /updates/, принимающий в теле запроса множество метрик в формате: []Metrics (списка метрик).
Агент:
  Научите агент работать с использованием нового API (отправлять метрики батчами).
  Стоит помнить, что:
  - нужно соблюдать обратную совместимость;
  - отправлять пустые батчи не нужно;
  - вы умеете сжимать контент по алгоритму gzip;
  - изменение в базе можно выполнять в рамках одной транзакции или одного запроса;
  - необходимо избегать формирования условий для возникновения состояния гонки (race condition).

# Инкремент 13
Измените весь свой код в соответствии со знаниями, полученными в этой теме. Добавьте обработку retriable-ошибок.
Retriable-ошибки — это ошибки, которые могут быть исправлены повторной попыткой выполнения операции. Это бывает полезно для программ, которые работают с сетью или файловой системой, где возможны временные проблемы связи или доступа к данным. Ошибки могут быть вызваны различными причинами, такими как перегрузка сервера, недоступность сети или ошибки в коде программы.
Примеры retriable-ошибок:
- Ошибка связи с сервером при отправке запроса.
- Ошибка чтения данных из сети или БД из-за проблем соединения.
- Ошибка доступа к файлу, который был заблокирован другим процессом.
Сценарии возможных ошибок:
- Агент не сумел с первой попытки выгрузить данные на сервер из-за временной невозможности установить соединение с сервером.
- При обращении к PostgreSQL cервер получил ошибку транспорта (из категории Class 08 — Connection Exception).
Стратегия реализации:
- Количество повторов должно быть ограничено тремя дополнительными попытками.
- Интервалы между повторами должны увеличиваться: 1s, 3s, 5s.
- Чтобы определить тип ошибки PostgreSQL, с которой завершился запрос, можно воспользоваться библиотекой github.com/jackc/pgerrcode, в частности pgerrcode.UniqueViolation.