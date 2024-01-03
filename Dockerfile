FROM golang:1.20

WORKDIR /btc-sbt

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN make build

ENTRYPOINT ["btc-sbt", "node"]
