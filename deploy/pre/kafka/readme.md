## 创建kafka topic
    默认是不允许程序自动创建topic的，进入kafka的容器创建3个topic
    进入容器
    $ docker exec -it kafka /bin/sh
    $ cd /opt/bitnami/kafka/bin/
    创建topic
    log专用kafka
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 10 --topic jakarta-log
    业务专用kafka
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic im-state-change-msg-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic update-payment-status-topic 
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic update-refund-status-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic update-order-action-topic 
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 1 --topic im-define-msg-send-topic 
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic update-chat-stat-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic check-chat-state-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic update-cash-status-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic commit-move-cash-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic subscribe-notify-msg-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic update-listener-user-stat-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic im-after-send-msg-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic wxfwh-callback-event-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic wx-mini-program-msg-send-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic wx-fwh-msg-send-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic upload-user-event-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic enter-chat-topic
    ./kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 -partitions 3 --topic rec-listener-when-user-login-topic