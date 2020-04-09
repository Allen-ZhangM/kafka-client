### kafka build

### java8

```
tar -zxvf jdk-8u211-linux-x64.tar.gz

环境变量：/etc/profile

JAVA_HOME=/home/admin/jdk1.8.0_211
PATH=${JAVA_HOME}/bin:$PATH

source /etc/profile

```

### zookeeper

kafka自带zookeeper，可以跳过这步

```
wget http://apache.osuosl.org/zookeeper/zookeeper-3.4.14/zookeeper-3.4.14.tar.gz
tar -zxvf zookeeper-3.4.14.tar.gz 

环境变量：/etc/profile

zkServer.sh start
zkServer.sh status

查看2181端口
netstat -tunple | grep 2181

参考链接：
https://blog.csdn.net/fengchen0123456789/article/details/82220522
```

### kafka

```
wget http://mirrors.tuna.tsinghua.edu.cn/apache/kafka/2.2.0/kafka_2.11-2.2.0.tgz
tar -zxvf kafka_2.11-2.2.0.tgz

启动zookeeper
bin/zookeeper-server-start.sh -daemon config/zookeeper.properties

启动broker
bin/kafka-server-start.sh -daemon config/server.properties

通过supervisor方式配置启动

查看9092端口
netstat -tunple | grep 9092

测试,kafka目录bin下
./kafka-console-producer.sh --broker-list 192.168.194.128:9092,192.168.194.128:9093,192.168.194.128:9094 --topic test
./kafka-console-consumer.sh --bootstrap-server 192.168.194.128:9092,192.168.194.128:9093,192.168.194.128:9094 --topic test

参考链接：
https://blog.csdn.net/qq_21057881/article/details/86763149
```

### zookeeper配置参数

```
dataDir=/tmp/zookeeper
clientPort=2181
host.name=localhost
```

### kafka server 配置参数

```
//id唯一
broker.id=0
//端口和ip
port=9092
host.name=192.168.194.128

//外网访问需要配置
listeners=PLAINTEXT://192.168.194.128:9092
advertised.listeners=PLAINTEXT://192.168.194.128:9092

//log目录
log.dirs=/tmp/kafka-logs
//分区数，相当于负载均衡，每个分区内消息有序，整体消息有序情况下设置为1.
num.partitions=1
//zk配置
zookeeper.connect=localhost:2181

//高可用集群，需要配置为多个
offsets.topic.replication.factor=1
transaction.state.log.replication.factor=1
transaction.state.log.min.isr=1
```

### 集群状态查询

```
查询topic
bin/kafka-topics.sh --describe --zookeeper localhost:2181 --topic __consumer_offsets
如果Leader均匀分布，Replicas和Isr有多个节点说明正常
Topic:__consumer_offsets	PartitionCount:50	ReplicationFactor:3	Configs:segment.bytes=104857600,cleanup.policy=compact,compression.type=producer
	Topic: __consumer_offsets	Partition: 0	Leader: 3	Replicas: 3,0,2	Isr: 3,0,2
	Topic: __consumer_offsets	Partition: 1	Leader: 0	Replicas: 0,2,3	Isr: 0,2,3
	Topic: __consumer_offsets	Partition: 2	Leader: 2	Replicas: 2,3,0	Isr: 2,3,0

```

### Kafka图形管理工具

http://www.kafkatool.com/download.html 

### 更新方案

```
目的：保证在更新过程中不会丢失消息，程序部署过程出现的异常有预知

1.在tracker端新增发送kafka消息，同样的消息用redis和kafka都发送一遍
2.在boss端新增接收kafka消息，只输出日志，对比redis和kafka消息数目、内容
3.把tracker端发送消息更改为kafka发送，新消息已经存在broker中了
4.等待redis队列为空，把boss端接收消息改为kafka
```

### kafka维护
```
//查看topic列表
bin/kafka-topics.sh --zookeeper 172.19.0.83:8300 --list
//查看group消费情况
bin/kafka-consumer-groups.sh --bootstrap-server 172.19.0.83:8310 --describe --group gc1
//删除topic
设置broker参数
delete.topic.enable = true
bin/kafka-topics.sh --delete --zookeeper 172.19.0.83:8300 --topic test
//查看有多少个组
bin/kafka-consumer-groups.sh --bootstrap-server 192.168.194.128:9092 --list
```

### 错误日志

```
Java HotSpot(TM) 64-Bit Server VM warning: INFO: os::commit_memory(0x00000000c0000000, 1073741824, 0) failed; error='Cannot allocate memory' (errno=12)
#
# There is insufficient memory for the Java Runtime Environment to continue.
# Native memory allocation (mmap) failed to map 1073741824 bytes for committing reserved memory.
# An error report file with more information is saved as:
# /home/admin/kafka_2.11-2.2.0/hs_err_pid16354.log

解决方案：
修改kafka-server-start.sh文件
找到这一行
export KAFKA_HEAP_OPTS="-Xmx1G -Xms1G" 
改为合适的大小
export KAFKA_HEAP_OPTS="-Xmx256M -Xms128M" 
```

```
[2019-07-11 15:18:10,639] ERROR [KafkaServer id=0] Fatal error during KafkaServer startup. Prepare to shutdown (kafka.server.KafkaServer)
org.apache.kafka.common.KafkaException: Socket server failed to bind to 139.224.119.126:9092: Cannot assign requested address.
	at kafka.network.Acceptor.openServerSocket(SocketServer.scala:573)
	at kafka.network.Acceptor.<init>(SocketServer.scala:451)
	at kafka.network.SocketServer.kafka$network$SocketServer$$createAcceptor(SocketServer.scala:245)
	at kafka.network.SocketServer$$anonfun$createDataPlaneAcceptorsAndProcessors$1.apply(SocketServer.scala:215)
	at kafka.network.SocketServer$$anonfun$createDataPlaneAcceptorsAndProcessors$1.apply(SocketServer.scala:214)
	at scala.collection.mutable.ResizableArray$class.foreach(ResizableArray.scala:59)
	at scala.collection.mutable.ArrayBuffer.foreach(ArrayBuffer.scala:48)
	at kafka.network.SocketServer.createDataPlaneAcceptorsAndProcessors(SocketServer.scala:214)
	at kafka.network.SocketServer.startup(SocketServer.scala:114)
	at kafka.server.KafkaServer.startup(KafkaServer.scala:253)
	at kafka.server.KafkaServerStartable.startup(KafkaServerStartable.scala:38)
	at kafka.Kafka$.main(Kafka.scala:75)
	at kafka.Kafka.main(Kafka.scala)
Caused by: java.net.BindException: Cannot assign requested address
	at sun.nio.ch.Net.bind0(Native Method)
	at sun.nio.ch.Net.bind(Net.java:433)
	at sun.nio.ch.Net.bind(Net.java:425)
	at sun.nio.ch.ServerSocketChannelImpl.bind(ServerSocketChannelImpl.java:223)
	at sun.nio.ch.ServerSocketAdaptor.bind(ServerSocketAdaptor.java:74)
	at sun.nio.ch.ServerSocketAdaptor.bind(ServerSocketAdaptor.java:67)
	at kafka.network.Acceptor.openServerSocket(SocketServer.scala:569)
	... 12 more
	
	
解决方案：
server.properties中地址配置有误
可能为以下3条地址有误
host.name=192.168.194.128
listeners=PLAINTEXT://192.168.194.128:9092
advertised.listeners=PLAINTEXT://192.168.194.128:9092
```



