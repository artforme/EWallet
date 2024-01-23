FROM golang:1.21.4

WORKDIR /go/src/EWallet

EXPOSE 8082

COPY . .

WORKDIR /go/src/EWallet/cmd/EWallet

ENV CONFIG_PATH /go/src/EWallet/config/local.yaml
ENV STORAGE_PATH /go/src/EWallet/storage/storage.db

RUN go mod download

RUN go build -o /go/src/EWallet/cmd/EWallet/main .

CMD ["/go/src/EWallet/cmd/EWallet/main"]