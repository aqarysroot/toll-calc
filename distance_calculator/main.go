package main

import (
	"fmt"
	"log"
	"toll-calculator/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000"
)

// Transport(HTTP, GRPC, kafka) -> attach business logic to this

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleWare(svc)

	httpClient := client.NewHTTPClient(aggregatorEndpoint)
	//grpcClient, err := client.NewGRPCClient()

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, httpClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()

	fmt.Println("Wokr fine")
}
