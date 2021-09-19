BUILD_DIR=build
CMD_DIR=cmd
CMDS=$(patsubst $(CMD_DIR)/%,%,$(wildcard $(CMD_DIR)/*))

.PHONY: fmt check test

all: fmt check test bin

fmt:
	gofmt -s -w -l .
	goimports -w $(shell find . -type f -name '*.go' -not -path "./internal/*")

check:
	golint -set_exit_status ./... && \
	go vet -all ./... && \
	misspell -error */** && \
	go mod tidy

test:
	go test ./...

bin:
	./build_docker.sh --bin

docker: $(CMDS)
	./build_docker.sh --docker $^

$(CMDS):
	./build_docker.sh --bin $@

proto:
		protoc --go_out=internal/rpc/hello --go_opt=module=protos --go-grpc_out=internal/rpc/hello --go-grpc_opt=module=protos protos/HelloService.proto

setup: setup-postgres

# add this in /etc/fstab and run `sudo mount -a`
# tmpfs /mtmp tmpfs size=2048m,mode=1777 0 0
setup-postgres:
	@if ! docker ps | /bin/grep postgres-localdev; then \
		docker run --name postgres-localdev \
			-p 5432:5432 \
			-d --tmpfs /var/lib/postgresql/data:rw,noexec,nosuid,size=4096m \
			-d --tmpfs /run:rw,noexec,nosuid,size=4096m \
			-e POSTGRES_DB=testing \
			-e POSTGRES_USER=tester \
			-e POSTGRES_PASSWORD=aaaa1234 \
			--restart always \
			-d postgres:12 \
			-c fsync=off -c full_page_writes=off; \
		docker run --rm --link postgres-localdev:postgres-localdev aanand/wait; \
	fi

remove:
	docker rm -f postgres-localdev
