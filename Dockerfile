FROM golang:1.21.6-alpine3.19
WORKDIR /app
COPY . /app
RUN go mod tidy
RUN make run
EXPOSE 3000
