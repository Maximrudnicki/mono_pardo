build:
	go build -o bin/mono_pardo.exe ./cmd/.

run: build
	./bin/mono_pardo

test:
	go test ./tests/... -count=1