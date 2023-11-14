gate:
	@go build -o bin/gate gateway/*.go
	@./bin/gate

obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

reciever:
	@go build -o bin/reciever data_reciever/*.go
	@./bin/reciever


calculator:
	@go build -o bin/calculator distance_calculator/*.go
	@./bin/calculator

agg:
	@go build -o bin/agg aggregator/*.go
	@./bin/agg

proto:
	protoc --go_out=. --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

.PHONY: obu agg