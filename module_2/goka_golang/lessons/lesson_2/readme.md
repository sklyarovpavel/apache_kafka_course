Блокировка и разблокировка пользователей
```bash
go run cmd/block-user/main.go -user user-0 -stream blocker-stream
go run cmd/block-user/main.go -user user-1 -stream blocker-stream
go run cmd/block-user/main.go -user user-2 -stream blocker-stream
```


```bash
go run cmd/block-user/main.go -unblock -user user-0 -stream blocker-stream
go run cmd/block-user/main.go -unblock -user user-1 -stream blocker-stream
go run cmd/block-user/main.go -unblock -user user-2 -stream blocker-stream
```


Создание топика для передачи данных от эммитера в фильтр
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic emitter2filter-stream --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2
```
Создание топика отслеживающего изменения в групповой таблицы лайков пользователей.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic user-like-table --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```


Создание топика для передачи данных от фильтра в обработчик пользовательских лайков
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic filter2userprocessor-stream --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2
```
Создание топика отслеживающего изменения в групповой таблицы фильтра.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic filter-table --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```



Создание топика для передачи данных заблокированных пользователей.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic blocker-stream --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```

Создание топика отслеживающего изменения в групповой таблицы заблокированных пользователей.
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic blocker-table --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```


Создание топика для передачи данных о статьях
```bash
docker exec -it kafka-1 ../../usr/bin/kafka-topics --create --topic namer-stream --bootstrap-server localhost:9092 --partitions 3 --replication-factor 2 --config cleanup.policy=compact
```


```bash
go run cmd/article-namer/main.go -word "1" -with "how to use goka" 
go run cmd/article-namer/main.go -word "0" -with "how to use kafka" 
go run cmd/article-namer/main.go -word "3" -with "how to use goka with kafka" 
```
