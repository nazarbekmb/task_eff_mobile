FROM golang:1.23.0

WORKDIR /api-server

COPY go.mod go.sum ./
RUN go mod download
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
COPY . ./

RUN go build -o main ./cmd/main.go

EXPOSE 8080

CMD ["./main"]