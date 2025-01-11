FROM golang:1.22.4-bullseye
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o geth-indexer
RUN chmod 777 ./geth-indexer
ENTRYPOINT [ "./geth-indexer" ]