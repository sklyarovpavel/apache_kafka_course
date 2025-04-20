# Проект модуля №3: "Администрирование кластера Kafka"

## Задание "Балансировка партиций и диагностика кластера"

**Цели задания:**
1. Освоить балансировку партиций и распределение нагрузки с 
помощью Partition Reassignment Tools.
2. Попрактиковаться в диагностике и устранении проблем кластера.

**Задание:**
1. Создайте новый топик `balanced_topic` с 8 партициями и фактором репликации 3.
2. Определите текущее распределение партиций.
3. Создайте JSON-файл `reassignment.json` для перераспределения партиций.
4. Перераспределите партиции.
5. Проверьте статус перераспределения.
6. Убедитесь, что конфигурация изменилась. 
7. Смоделируйте сбой брокера:
   * Остановите брокер `kafka-1`.
   * Проверьте состояние топиков после сбоя.
   * Запустите брокер заново. 
   * Проверьте, восстановилась ли синхронизация реплик.

## Решение

1. **Создайте новый топик `balanced_topic` с 8 партициями и фактором 
   репликации 3.**

    Запускаем Docker compose (локальный терминал):
    ```bash
    docker compose up -d
    ```
    
    Переходим в брокера `kafka-0` и выполняем команды в терминале:
    
    a. Проверяем, что нет созданных топиков (kafka-0 терминал):
    ```bash
    kafka-topics --bootstrap-server localhost:9092 --list
    ```
    Должны увидеть пустой вывод.

    b. Создаем новый топик `balanced_topic` с 8 партициями и фактором 
       репликации 3 (kafka-0 терминал):
    ```bash
    kafka-topics --bootstrap-server localhost:9092 --create --topic balanced_topic --replication-factor 3 --partitions 8
    ```
    Должны увидеть `Created topic balanced_topic.`.

    c. Проверяем, что топик создался (kafka-0 терминал):
    ```bash
    kafka-topics --bootstrap-server localhost:9092 --list
    ```
    Должны увидеть новый топик `balanced_topic`.

2. **Определите текущее распределение партиций.**

    Выполним команду для просмотра конфигурации топика `balanced_topic` 
    (kafka-0 терминал):
    ```bash
    kafka-topics --bootstrap-server localhost:9092 --topic balanced_topic --describe
    ```
    Должны увидеть 8 партиций и фактор репликации 3. Пример вывода:
    ```
    Topic: balanced_topic   TopicId: JPAICzpbTFC147-_Girw_g PartitionCount: 8       ReplicationFactor: 3    Configs: 
        Topic: balanced_topic   Partition: 0    Leader: 2       Replicas: 2,0,1 Isr: 2,0,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 1    Leader: 0       Replicas: 0,1,2 Isr: 0,1,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 2    Leader: 1       Replicas: 1,2,0 Isr: 1,2,0      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 3    Leader: 0       Replicas: 0,1,2 Isr: 0,1,2      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 4    Leader: 1       Replicas: 1,2,0 Isr: 1,2,0      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 5    Leader: 2       Replicas: 2,0,1 Isr: 2,0,1      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 6    Leader: 1       Replicas: 1,2,0 Isr: 1,2,0      Elr:    LastKnownElr: 
        Topic: balanced_topic   Partition: 7    Leader: 2       Replicas: 2,0,1 Isr: 2,0,1      Elr:    LastKnownElr: 
    ```

3. **Создайте JSON-файл `reassignment.json` для перераспределения партиций.**

    Создаем JSON-файл `topics_to_move.json` для генерации перераспределения 
    партиций (в контейнере kafka-0):
    ```json
    {
        "version": 1,
        "topics": [{"topic": "balanced_topic"}]
    }
    ```
    
    Генерируем план нового распределения (kafka-0 терминал):
    ```bash
    kafka-reassign-partitions --bootstrap-server localhost:9092 --broker-list "0,1,2" --topics-to-move-json-file "topics_to_move.json" --generate
    ```
    Должны получить что-то вроде:
    ```
    Current partition replica assignment
    {"version":1,"partitions":[{"topic":"balanced_topic","partition":0,"replicas":[2,0,1],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":1,"replicas":[0,1,2],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":2,"replicas":[1,2,0],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":3,"replicas":[0,1,2],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":4,"replicas":[1,2,0],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":5,"replicas":[2,0,1],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":6,"replicas":[1,2,0],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":7,"replicas":[2,0,1],"log_dirs":["any","any","any"]}]}

    Proposed partition reassignment configuration
    {"version":1,"partitions":[{"topic":"balanced_topic","partition":0,"replicas":[1,2,0],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":1,"replicas":[2,0,1],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":2,"replicas":[0,1,2],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":3,"replicas":[1,0,2],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":4,"replicas":[2,1,0],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":5,"replicas":[0,2,1],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":6,"replicas":[1,2,0],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":7,"replicas":[2,0,1],"log_dirs":["any","any","any"]}]}
    ```
   
    Копируем вывод из `Proposed partition reassignment configuration` и 
    записываем его в файл `reassignment.json`.

    > p.s. Для записи в файл можно использовать `echo "any_text" >> file.json` 
    или `nano`/`vim`/etc. 
    > 
    > Установка vim в контейнер:
    > ```bash
    > docker exec -it --user root kafka-0 /bin/bash
    > yum install vim
    > ```

4. **Перераспределите партиции.**

    Выполняем план перераспределения партиций (kafka-0 терминал):
    ```bash
    kafka-reassign-partitions --bootstrap-server localhost:9092 --reassignment-json-file reassignment.json --execute
    ```
    Должны увидеть что-то подобное:
    ```
    Current partition replica assignment

    {"version":1,"partitions":[{"topic":"balanced_topic","partition":0,"replicas":[2,0,1],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":1,"replicas":[0,1,2],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":2,"replicas":[1,2,0],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":3,"replicas":[0,1,2],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":4,"replicas":[1,2,0],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":5,"replicas":[2,0,1],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":6,"replicas":[1,2,0],"log_dirs":["any","any","any"]},{"topic":"balanced_topic","partition":7,"replicas":[2,0,1],"log_dirs":["any","any","any"]}]}

    Save this to use as the --reassignment-json-file option during rollback
    Successfully started partition reassignments for balanced_topic-0,balanced_topic-1,balanced_topic-2,balanced_topic-3,balanced_topic-4,balanced_topic-5,balanced_topic-6,balanced_topic-7
    ```

5. **Проверьте статус перераспределения.**

    Выполним команду для проверки статуса перераспределения партиций 
    (kafka-0 терминал):
    ```bash
    kafka-reassign-partitions --bootstrap-server localhost:9092 --reassignment-json-file reassignment.json --verify
    ```
    После успешного перераспределения вывод должен выглядеть примерно 
    так:
    ```
    Status of partition reassignment:
    Reassignment of partition balanced_topic-0 is completed.
    Reassignment of partition balanced_topic-1 is completed.
    Reassignment of partition balanced_topic-2 is completed.
    Reassignment of partition balanced_topic-3 is completed.
    Reassignment of partition balanced_topic-4 is completed.
    Reassignment of partition balanced_topic-5 is completed.
    Reassignment of partition balanced_topic-6 is completed.
    Reassignment of partition balanced_topic-7 is completed.

    Clearing broker-level throttles on brokers 0,1,2
    Clearing topic-level throttles on topic balanced_topic
    ```

6. **Убедитесь, что конфигурация изменилась.**

    Выполним команду для просмотра конфигурации топика `balanced_topic` 
    (kafka-0 терминал):
    ```bash
    kafka-topics --bootstrap-server localhost:9092 --topic balanced_topic --describe
    ```
    Должны увидеть **новое** распределение (которое указано в файле 
    `reassignment.json`).

7. **Смоделируйте сбой брокера:**
   * Остановите брокер `kafka-1`.

     Выполните команду остановки контейнера (локальный терминал):
     ```bash
     docker stop kafka-1
     ```

   * Проверьте состояние топиков после сбоя.
   
     Выполним команду для просмотра конфигурации топика `balanced_topic` 
     (kafka-0 терминал):
     ```bash
     kafka-topics --bootstrap-server localhost:9092 --topic balanced_topic --describe
     ```
     
     В выводе в полях `Isr:` пропал брокер kafka-1:
     ```
     Topic: balanced_topic   TopicId: JPAICzpbTFC147-_Girw_g PartitionCount: 8       ReplicationFactor: 3    Configs: 
            Topic: balanced_topic   Partition: 0    Leader: 2       Replicas: 1,2,0 Isr: 2,0        Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 1    Leader: 2       Replicas: 2,0,1 Isr: 0,2        Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 2    Leader: 0       Replicas: 0,1,2 Isr: 2,0        Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 3    Leader: 0       Replicas: 1,0,2 Isr: 0,2        Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 4    Leader: 2       Replicas: 2,1,0 Isr: 2,0        Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 5    Leader: 0       Replicas: 0,2,1 Isr: 2,0        Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 6    Leader: 2       Replicas: 1,2,0 Isr: 2,0        Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 7    Leader: 2       Replicas: 2,0,1 Isr: 2,0        Elr:    LastKnownElr: 
     ```
     
   * Запустите брокер заново.
   
     Выполните команду запуска контейнера (локальный терминал):
     ```bash
     docker start kafka-1
     ```

   * Проверьте, восстановилась ли синхронизация реплик.

     Выполним команду для просмотра конфигурации топика `balanced_topic` 
     (kafka-0 терминал):
     ```bash
     kafka-topics --bootstrap-server localhost:9092 --topic balanced_topic --describe
     ```
     
     В выводе в полях `Isr:` вернулся брокер kafka-1:
     ```
     Topic: balanced_topic   TopicId: JPAICzpbTFC147-_Girw_g PartitionCount: 8       ReplicationFactor: 3    Configs: 
            Topic: balanced_topic   Partition: 0    Leader: 2       Replicas: 1,2,0 Isr: 2,0,1      Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 1    Leader: 2       Replicas: 2,0,1 Isr: 0,2,1      Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 2    Leader: 0       Replicas: 0,1,2 Isr: 2,0,1      Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 3    Leader: 0       Replicas: 1,0,2 Isr: 0,2,1      Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 4    Leader: 2       Replicas: 2,1,0 Isr: 2,0,1      Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 5    Leader: 0       Replicas: 0,2,1 Isr: 2,0,1      Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 6    Leader: 2       Replicas: 1,2,0 Isr: 2,0,1      Elr:    LastKnownElr: 
            Topic: balanced_topic   Partition: 7    Leader: 2       Replicas: 2,0,1 Isr: 2,0,1      Elr:    LastKnownElr: 
     ```
     Вывод: синхронизация восстановлена.