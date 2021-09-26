package controller

import (
	"analytix/generated"
	"database/sql"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"time"
)

type EventProcessor interface {
	ProcessEventBatch()
}

type KafkaEventProcessor struct {
	kafkaConsumer        *kafka.Consumer
	clickhouseConnection *sql.DB
}

func (k KafkaEventProcessor) ProcessEventBatch() {
	// Process events batch.
	run := true
	msgCount := 0
	var startTime = time.Now()
	var buffer = make(chan []*generated.AcquisitionEvent)

	var localBatch []*generated.AcquisitionEvent
	var batchSize = 100
	var batchTimeMs int64 = 300
	var currentBatchCount = 0
	go func() {
		for {
			if localBatch == nil {
				localBatch = make([]*generated.AcquisitionEvent, 10)
			}
			currentBatchCount++
			ev := k.kafkaConsumer.Poll(0)
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				k.kafkaConsumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				k.kafkaConsumer.Unassign()
			case *kafka.Message:
				msgCount += 1
				msg := ev.(*kafka.Message)
				var event generated.AcquisitionEvent
				err := proto.Unmarshal(msg.Value, &event)
				if err != nil {
					println("Unable to process message " + err.Error())
				}

				localBatch = append(localBatch, &event)
				if time.Now().Sub(startTime).Milliseconds() > batchTimeMs || currentBatchCount > batchSize {
					println("Flushing batch.")
					buffer <- localBatch
					localBatch = nil
				}
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			default:
				//fmt.Printf("Ignored %v\n", e)
			}
		}
	}()

	for batch := range buffer {
		// JDBC Send to Clickhouse.
		count, err := k.StreamToClickhouse(batch)
		if ( err != nil){
			println(err.Error())
		}
		fmt.Printf("Inserted %d records into analytics Datastore. \n",count)

	}
}

func (k *KafkaEventProcessor) StreamToClickhouse(batch []*generated.AcquisitionEvent) (int32,error) {
	if batch == nil {
		return 0 ,nil
	}
	// Extract user id,

	tx, err := k.clickhouseConnection.Begin()

	if ( err != nil){
		return 0,nil
	}

	stmt, err := tx.Prepare(`
		"INSERT INTO example (
			event_ts, 
			client_id, 
			event_key, 
			guest_id, 
			user_id, 
			session_id,
			event_name, 
			event_type,
			channel,
			data) VALUES (?, ?, ?, ?, ?, ? , ?)")
	`);

	defer stmt.Close()
	var eventCount int32 = 0
	for i := 0; i < len(batch); i++ {
		currentBatchEvents := batch[i]
		clientId := currentBatchEvents.ClientId
		for _, payload := range currentBatchEvents.Payload {
			var args []interface{}
			guestId := 				payload.GuestId
			userId := 				payload.UserId
			sessionId := 				payload.SessionId
			eventName := 				payload.EventName
			eventType := 				payload.EventType
			channel := 				payload.Channel
			args = append(args, payload.EventTsNanos, clientId,guestId,userId,sessionId,eventName,eventType,channel)
			getMapKeyValuesFlattened(payload.Data,&args)

			if _, err := stmt.Exec( args); err != nil {
				log.Fatal(err)
			}
			eventCount++;
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return eventCount,nil

}

func getMapKeyValuesFlattened(dataMap map[string]*structpb.Value,existingArgs *[]interface{})  {
	for key , mapValue := range dataMap {
		stringValue := mapValue.String()
		*existingArgs = append(*existingArgs, key,stringValue)
	}
}

func NewEventsProcessorConsumer(kConsumer *kafka.Consumer, clickhouseConnection *sql.DB) EventProcessor {
	ac := KafkaEventProcessor{
		kafkaConsumer:        kConsumer,
		clickhouseConnection: clickhouseConnection,
	}
	return ac
}
