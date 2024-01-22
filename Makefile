run-tests:
	go test ./...
	
deploy:
	docker compose -f ./deployment/compose.yaml up  