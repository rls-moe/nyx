FROM golang:1.14-buster AS build-env

COPY . ${GOPATH}/src/github.com/rls-moe/nyx
WORKDIR ${GOPATH}/src/github.com/rls-moe/nyx

RUN go build -o $GOPATH/bin/nyx

FROM debian:bullseye
LABEL maintainer="b.pedini@bjphoster.com"

EXPOSE 8080

RUN groupadd \
    --gid 1000 \
    nyx && \
    useradd \
    --home-dir /opt/nyx \
    --comment "Nyx Board" \
    --gid nyx \
    --create-home \
    --no-user-group \
    --uid 1000 \
    --shell /bin/bash \
    nyx

COPY --from=build-env /go/bin/nyx /opt/nyx/nyx
COPY --from=build-env /go/src/github.com/rls-moe/nyx/config.example.yml /opt/nyx/config.yml
RUN chown -R \
    nyx:nyx /opt/nyx

WORKDIR /opt/nyx
USER nyx
CMD [ "/opt/nyx/nyx" ]
