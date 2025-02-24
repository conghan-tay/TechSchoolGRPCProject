gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb --grpc-gateway_out=:pb  --openapiv2_out=:openapi

clean:
	rm pb/*.go

server:
	go run cmd/server/main.go -port 8800

client:
	go run cmd/client/main.go -address 0.0.0.0:8800

test:
	go test -cover -race ./...

cert:
	cd cert; ./gen.sh; cd ..

.PHONY: gen clean server client test cert