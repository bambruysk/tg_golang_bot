FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go mod vendor
RUN go build -o bot

CMD ["bot"]

#ENTRYPOINT bot
