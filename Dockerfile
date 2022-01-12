FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY server.go ./
COPY gqlgen.yml ./
COPY graph ./graph
COPY models ./models

RUN go build -o /docker-gql-server ./server.go

EXPOSE 8080

CMD [ "/docker-gql-server" ]