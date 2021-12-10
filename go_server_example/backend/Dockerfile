FROM golang:1.17

RUN apt-get update \
    && apt-get install -y --no-install-recommends postgresql-client \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/backend
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["backend"]
