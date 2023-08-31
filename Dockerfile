FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o go-app ./cmd/main.go
RUN go test ./internal/http-server/handlers/createSeg ./internal/http-server/handlers/deleteSeg ./internal/http-server/handlers/changeUser ./internal/http-server/handlers/getClientSeg
CMD ["./go-app"]