package db

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"

)

type KafkaConsumerConfig struct {
	BootstrapServers string `envconfig:"KAFKA_BOOTSTRAP_SERVERS" required:"true" default:"localhost:9092"`
	GroupId string `envconfig:"KAFKA_CONSUMER_GROUP_ID" default:"sample-consumer"`
	OffsetReset string `envconfig:"KAFKA_CONSUMER_OFFSET_RESET"`
}

type KafkaConfig struct {
	BootstrapServers string `envconfig:"KAFKA_BOOTSTRAP_SERVERS" required:"true" default:"localhost:9092"`
	TcpPort string `envconfig:"KAFKA_CLIENT_ID_PREFIX"`
	PingOnConnect bool `envconfig:"KAFKA_DB_PING_ON_CONNECT"`
}

func NewKafkaProducer(config *KafkaConfig)(*kafka.Producer,error){
	hostname , err := os.Hostname() ;
	if ( err != nil){
		return nil,err;
	}
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.BootstrapServers,
		"client.id": hostname,
		"acks": "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		return nil, err
	}
	return p , nil
}


func NewKafkaConsumer(config *KafkaConsumerConfig) (*kafka.Consumer, error){
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    config.BootstrapServers,
		"group.id":             config.GroupId,
		"auto.offset.reset":   config.OffsetReset})
	if ( err != nil){
		return nil,err
	}
	err = consumer.SubscribeTopics([]string{"acquisition_event_stream"}, nil)
	if err != nil {
		return nil, err
	}
	return consumer, nil;
}