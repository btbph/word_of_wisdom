run-tests:
	go test ./...
	
run-server:
	go run ./cmd/server/server.go

run-client:
	go run ./cmd/client/client.go

deploy:
	docker compose -f ./deployment/compose.yaml up  