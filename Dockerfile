FROM golang:alpine 

WORKDIR /app
COPY . /app
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
RUN go build -o ./bin/qtun ./main.go

ENTRYPOINT ["./bin/qtun"]

