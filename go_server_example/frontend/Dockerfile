FROM golang:1.17

WORKDIR /go/src/frontend
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["frontend"]
