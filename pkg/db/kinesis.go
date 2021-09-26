package db

import (
	"analytix/pkg/lib"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

type KinesisConfig struct {
	BootstrapServers string `envconfig:"KAFKA_BOOTSTRAP_SERVERS" required:"true" default:"localhost:9092"`
	TcpPort string `envconfig:"KAFKA_CLIENT_ID_PREFIX"`
	PingOnConnect bool `envconfig:"KAFKA_DB_PING_ON_CONNECT"`
}

func NewKinesisClient()(*kinesis.Client){
	config := lib.NewAwsConfigForService(lib.ServiceKinesis)
	kinesisClient  := kinesis.NewFromConfig(config)
		return kinesisClient
}