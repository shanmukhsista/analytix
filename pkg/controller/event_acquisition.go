package controller

import (
	"analytix/generated"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
	"net/http"
)

type AcquisitionServiceKafkaController struct {
	kProducer * kafka.Producer
	// Which channel to send eventss to.
	deliveryChannel chan kafka.Event
	eventAcquisitionTopic string
}

func (a AcquisitionServiceKafkaController) AcquireEvents(eventPayload *generated.AcquisitionEvent) (*generated.AcquisitionEventResponse, error) {
	// Publish message to kafka. Serialize to protobuf..
	// Assign a batch id.
	serializedValue , err:= proto.Marshal(eventPayload)
	if ( err != nil){
		return nil,err
	}
	err = a.kProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &a.eventAcquisitionTopic, Partition: kafka.PartitionAny},
		Value: serializedValue},
		a.deliveryChannel,
	)
	if ( err != nil){
		return nil, err
	}
	response := generated.AcquisitionEventResponse{
		Status:         http.StatusOK,
		ProcessingTime: 0,
		BatchId:        "",
	}
	return &response, nil
}

func NewAcquisitionServiceControllerKafka(kProducer * kafka.Producer, bufferSize int) EventAcquisitionAPI {
	ac := AcquisitionServiceKafkaController{
		kProducer:             kProducer,
		deliveryChannel:       make(chan kafka.Event,bufferSize),
			eventAcquisitionTopic: "acquisition_event_stream",
	}
	return ac
}

