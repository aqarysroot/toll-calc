package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"time"
	"toll-calculator/aggregator/client"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "the listen address")
	flag.Parse()
	aggregatorServiceAddr := flag.String("aggServiceAddr", "http://localhost:3000", "the listen address")

	var (
		client     = client.NewHTTPClient(*aggregatorServiceAddr)
		invHandler = newInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice))
	logrus.Info("gateway HTTP server running on port 6000")
	err := http.ListenAndServe(*listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

type InvoiceHandler struct {
	client client.Client
}

func newInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	// access agg client
	inv, err := h.client.GetInvoice(context.Background(), 43)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
			}).Info("REQ :: ")
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})

		}
	}
}
