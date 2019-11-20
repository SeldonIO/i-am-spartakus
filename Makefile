BINARY:=spartakus-linux-amd64
build_binary:
	mkdir -p ./bin/
	cd spartakus && \
		go build -o ../bin/$(BINARY) ./cmd/spartakus/main.go ./cmd/spartakus/collector.go ./cmd/spartakus/volunteer.go
	./bin/$(BINARY) collector --version

clean:
	rm -rfv ./bin/

