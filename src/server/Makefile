FLAGS=CGO_ENABLED=0 GOOS=linux GOARCH=arm

ENTRY=cmd/server/main.go

build-arm:
	$(FLAGS) go build -o main_arm $(ENTRY)

build:
	go build -o main $(ENTRY) 

clean: rm -f main
