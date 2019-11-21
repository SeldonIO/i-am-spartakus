IMAGE_NAME:=docker.io/seldonio/spartakus-amd64
IMAGE_TAG:=v1.1.0-hack-0.2

PKG:=github.com/kubernetes-incubator/spartakus
BINARY:=spartakus-linux-amd64

build_image:
	docker build --file=Dockerfile --force-rm=true -t $(IMAGE_NAME):$(IMAGE_TAG) .

push_image:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)

check_image_repo:
	xdg-open https://cloud.docker.com/u/seldonio/repository/docker/seldonio/spartakus-amd64


build_binary: clean
	mkdir -p ./bin/
	cd spartakus && \
		go build -o ../bin/$(BINARY) \
			-ldflags "-X $(PKG)/pkg/version.VERSION=$(IMAGE_TAG)" \
			./cmd/spartakus/main.go ./cmd/spartakus/collector.go ./cmd/spartakus/volunteer.go
	./bin/$(BINARY) collector --version

clean:
	rm -rfv ./bin/

