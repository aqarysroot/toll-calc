package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"toll-calculator/aggregator/client"
	"toll-calculator/types"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	store := NewMemoryStore()
	var (
		svc            = NewInvoiceAggregator(store)
		grpcListedAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		httpListedAddr = os.Getenv("HTTP_GRPC_ENDPOINT")
	)
	//svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(grpcListedAddr, svc))
	}()
	time.Sleep(time.Second * 2)
	c, err := client.NewGRPCClient(grpcListedAddr)
	if err != nil {
		log.Fatal(err)
	}
	err = c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 58.55,
		Unix:  time.Now().UnixNano(),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(makeHTTPTransport(httpListedAddr, svc))

}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	// make a TCP listener
	fmt.Println("GRPC Running on a port", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	// Make a new GRPC native server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	types.RegisterAggregatorServer(server, NewAggregatorGRPCServer(svc))
	// Register (OUR) GRPC server
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	//aggMetricHandler := newHTTPMetricsHandler("aggregate")
	invMetricHandler := newHTTPMetricsHandler("invoice")
	//http.HandleFunc("/aggregate", aggMetricHandler.instrument(handleAggregate(svc)))
	http.HandleFunc("/invoice", invMetricHandler.instrument(handleGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("server Running on a port", listenAddr)
	return http.ListenAndServe(listenAddr, nil)
}

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type given %s", storeType)
		return nil
	}
}

func writeJSON(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}
