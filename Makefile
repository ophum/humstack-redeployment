
run-apiserver:
	go run cmd/apiserver/main.go

run-agent:
	sudo go run cmd/agent/main.go --config cmd/agent/config.yaml

.PHONY: run-apiserver run-agent