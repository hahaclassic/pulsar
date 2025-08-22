build:
	go build -o pulsar ./cmd/pulsar/main.go

run:
	go run ./cmd/pulsar/main.go

stress-test: 
	go run ./tools/cpu_stress.go -n=12
