BIN=kubectl-guestfs 
build:
	mkdir -p bin/
	go build -o bin/$(BIN) main.go

image:
	docker build -t libguestfs-tools -f dockerfiles/Dockerfile .

install:
	install bin/$(BIN) /usr/local/bin

clean:
	rm -rf bin/
