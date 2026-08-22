FROM golang:1.21-alpine
ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=0
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o /app/bin/im-equiv .
CMD ["/app/bin/im-equiv", "help"]
