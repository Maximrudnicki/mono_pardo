build:
	go build -o bin/mono_pardo.exe ./cmd/.

run: build
	./bin/mono_pardo

fmt:
	go fmt ./...

test:
	go test ./tests/... -count=1

profile:
	go test ./tests/... -count=1 -cpuprofile cpu.prof -memprofile mem.prof