FROM golang:alpine as builder
WORKDIR /src
COPY . .
RUN apk add --no-cache --update git && \
    CGO_ENABLED=0 go build -o apiserver


FROM alpine:3

EXPOSE 53/udp 8081/tcp

ENV BASE_DOMAIN=
ENV DATA_PATH=/var/bind/zones

RUN apk add --update --no-cache curl
RUN set -ex; \
        apkArch="$(apk --print-arch)"; \
        case "$apkArch" in \
                armhf) arch='armhf' ;; \
                armv7) arch='arm' ;; \
                aarch64) arch='aarch64' ;; \
                x86_64) arch='amd64' ;; \
                *) echo >&2 "error: unsupported architecture: $apkArch"; exit 1 ;; \
        esac; \
                curl -L -o- https://github.com/just-containers/s6-overlay/releases/download/v2.2.0.3/s6-overlay-$arch.tar.gz | \
                 tar xfz - -C /

RUN addgroup -g 1000 named && adduser -h /etc/bind -s /sbin/nologin -G named -D -u 1000 named && \
    apk add --update --no-cache bind

COPY --from=builder /src/apiserver /apiserver
COPY rootfs /

VOLUME ${DATA_PATH}
WORKDIR ${DATA_PATH}
ENTRYPOINT ["/init"]
