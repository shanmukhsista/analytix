package main

import (
	"analytix/pkg/controller"
	"analytix/pkg/db"
	"analytix/pkg/lib"
	"analytix/pkg/service"
	"github.com/kelseyhightower/envconfig"
	"log"
)

// Dependencies Init here .
var acquisitionController controller.EventAcquisitionAPI

func init()  {
	var kafkaConfig db.KafkaConfig
	err := envconfig.Process( "",&kafkaConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	kafkaClient  , err :=  db.NewKafkaProducer(&kafkaConfig)
	if ( err != nil){
		panic(err)
	}
	acquisitionController  = controller.NewAcquisitionServiceControllerKafka(kafkaClient,1000000)
	if ( err != nil){
		log.Fatal(err.Error())
	}
}

func main(){

	var serverConfig lib.RestServerConfig
	err := envconfig.Process( "",&serverConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize Dependencies.
	acquisitionService := service.NewRestAcquisitionService(acquisitionController)
	server := lib.NewRestServer(&serverConfig)
	eventsGroup := server.GetEngine().Group("/api/v1/events")
	{
		eventsGroup.POST("/",acquisitionService.GetEvents)
	}
	err = server.Start()
	if ( err != nil){
		log.Fatal(err.Error())
	}
}
