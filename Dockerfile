FROM ubuntu:18.04

RUN apt-get update && \
    apt-get install ca-certificates -y
RUN update-ca-certificates

ADD bin/spartakus-linux-amd64 /spartakus

USER nobody:nobody
ENTRYPOINT ["/spartakus"]

