FROM golang:1.15 AS build

WORKDIR /src
# enable modules caching in separate layer
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -trimpath -ldflags "-s -w" -o dist/ethexporter ./ethexporter

FROM debian:10.2-slim AS runtime

ENV PORT 9890
ENV DATA /app

RUN mkdir -p /app && chown nobody:nogroup /app

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y \
        ca-certificates; \
    apt-get clean; \
    rm -rf /var/lib/apt/lists/*;

COPY --from=build /src/dist/ethexporter /usr/local/bin/ethexporter
COPY data/addresses.txt /app/addresses.txt

USER nobody
VOLUME /app
WORKDIR /app

ENTRYPOINT ["ethexporter"]
