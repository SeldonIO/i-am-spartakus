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

ADD bin/spartakus-linux-amd64 /spartakus

USER nobody:nogroup
ENTRYPOINT ["/spartakus"]

