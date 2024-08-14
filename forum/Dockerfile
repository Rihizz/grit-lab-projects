FROM golang:1.19-alpine
RUN apk add sqlite && apk add gcc && apk add g++
WORKDIR /forum
COPY go.mod ./
RUN go mod download
COPY ./ .
RUN go build -o /docker
EXPOSE 80:80
CMD [ "/docker" ]