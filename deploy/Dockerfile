FROM golang:alpine
WORKDIR /app
COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN go get -d -v ./...
RUN go build -o /project -v ./cmd/project/main.go
EXPOSE 8080
ENV CONN_STRING="test"
CMD ["/project"]