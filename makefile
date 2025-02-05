build:
	@go build -o ./bin/PlayLog
run: build
	@./bin/PlayLog
