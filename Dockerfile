FROM golang:1.21.6-alpine3.19
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build cmd/graphqlserver/main.go
CMD [ "./main" ]
EXPOSE 3000
