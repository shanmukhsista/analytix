package main

import (
	"analytix/pkg/controller"
	"analytix/pkg/db"
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Dependencies Init here .
var eventProcessor controller.EventProcessor

func init()  {
	var kafkaConfig db.KafkaConsumerConfig
	err := envconfig.Process( "",&kafkaConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	kafkaClient  , err :=  db.NewKafkaConsumer(&kafkaConfig)
	if ( err != nil){
		panic(err)
	}
	eventProcessor  = controller.NewEventsProcessorConsumer(kafkaClient)
	if ( err != nil){
		log.Fatal(err.Error())
	}
}

func main(){
	eventProcessor.ProcessEventBatch()
}
