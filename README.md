
# Introduction 

Service for Publishing events to Kafka.


## Programming Model.

## cmd

Any Applications that may be deployed to cloud.. All programs within cmd directory
are mostly wrappers around pkg services.

## pkg/*

All Core library and packages. 

### service

Any service that is exposed to users will belong here.

### Lib 

Core Library. Must always contain stateless logic that can be used by any module.

### Controller

This is the middle layer that helps APIs talk to the Database. Any interfacing / conversions must 
happen in this layer. 4

## Table Schema 

# Raw Events 

```sql
create table raw_events (
                            event_ts DateTime ,
                            client_id String ,
                            event_key String,
                            guest_id String ,
                            user_id String ,
                            session_id String,
                            event_name String,
                            data Map(String, String)
)ENGINE MergeTree()
partition by ( client_id,toYYYYMMDD(event_ts) )  order by (client_id,event_ts, user_id , guest_id);

``` 

