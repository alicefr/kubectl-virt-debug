BIN=kubectl-guestfs 
BIN_INIT_CONT=check-pvc
build:
	mkdir -p bin/
	go build -o bin/$(BIN) main.go
	go build -o bin/$(BIN_INIT_CONT) cmd/initcontainer/main.go

images: image-libguestfs image-init

image-libguestfs: build
	docker build -t libguestfs-tools -f dockerfiles/Dockerfile .

image-init: build
	docker build -t check-pvc -f dockerfiles/init/Dockerfile .

install:
	install bin/$(BIN) /usr/local/bin

clean:
	rm -rf bin/
