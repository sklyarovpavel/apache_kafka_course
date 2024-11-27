
Чтобы создать стрим из топика user-like-group-table, который будет содержать информацию о лайках и дизлайках для каждой статьи, вы можете использовать следующий запрос:

CREATE OR REPLACE STREAM user_like_ksqldb_stream (
"PostLike" MAP<STRING, BOOLEAN>
) WITH (
KAFKA_TOPIC = 'user-like-group-table',
VALUE_FORMAT = 'JSON'
);





CREATE OR REPLACE STREAM user_like_ksqldb_stream_with_timestamp1 AS
SELECT
"PostLike",
ROWTIME AS "EventTime"  -- Используем ROWTIME для получения временной метки
FROM user_like_ksqldb_stream EMIT CHANGES;



Пример KSQL запроса тут сделаем стрим как в гоке и сделаем такую же обработку:
CREATE OR REPLACE TABLE user_like_table (
"Key" STRING PRIMARY KEY,  -- Указываем поле Key как первичный ключ
"PostLike" MAP<STRING, BOOLEAN>
) WITH (
KAFKA_TOPIC = 'user-like-group-table',
VALUE_FORMAT = 'JSON'
);



CREATE TABLE queryable_user_like_table1 AS
SELECT *
FROM user_like_table;





SELECT * FROM queryable_user_like_table1;




Вывод для 1 результата

SELECT *
FROM queryable_user_like_table1
WHERE TRIM("Key") = 'user-1';


Или с помощью restapi


```bash
curl -X POST \
http://localhost:8088/query \
-H "Content-Type: application/json; charset=utf-8" \
-d '{
 "ksql": "SELECT * FROM queryable_user_like_table;",
 "streamsProperties": {}
}'
```



