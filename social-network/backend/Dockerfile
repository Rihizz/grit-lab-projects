FROM golang:1.19-alpine
RUN apk add sqlite && apk add gcc && apk add g++
WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o main .

EXPOSE 6969

CMD ["./main"]
