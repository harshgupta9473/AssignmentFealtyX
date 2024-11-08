FROM golang:1.22-alpine3.19
WORKDIR /test
COPY  . /test
RUN go mod tidy
RUN go build -o myapp .
EXPOSE 8000
ENTRYPOINT [ "./myapp" ]