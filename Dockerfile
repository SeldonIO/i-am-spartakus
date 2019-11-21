FROM golang:1.12 as builder

WORKDIR /work

# Download dependencies first to optimize Docker caching.
RUN mkdir -p spartakus
COPY spartakus/go.mod /work/spartakus/
COPY spartakus/go.sum /work/spartakus/
RUN cd spartakus && go mod download

COPY . .

RUN make build_binary


FROM ubuntu:18.04

RUN set -ex \
    && apt-get update -yqq \
    && apt-get install -yqq --no-install-recommends \
    ca-certificates \
    && apt-get clean \
    && rm -rf \
    /var/lib/apt/lists/* \
    /tmp/* \
    /var/tmp/* \
    /usr/share/man \
    /usr/share/doc \
    /usr/share/doc-base

RUN update-ca-certificates

COPY --from=builder /work/bin/spartakus-linux-amd64 /spartakus

USER nobody:nogroup
ENTRYPOINT ["/spartakus"]

