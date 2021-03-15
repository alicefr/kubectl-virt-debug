build:
	mkdir -p bin/
	go build -o bin/kubectl-virt-guestfs main.go

image:
	docker build -t libguestfs-tools -f dockerfiles/Dockerfile .


