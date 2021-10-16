#!bin/bash
export GO111MODULE=on
export GOPROXY=https://goproxy.cn

#Linux amd64
GOOS=linux GOARCH=amd64 go build -o ./bin/qtun-linux-amd64 ./main.go
#Linux arm64
GOOS=linux GOARCH=arm64 go build -o ./bin/qtun-linux-arm64 ./main.go
#Mac amd64
GOOS=darwin GOARCH=amd64 go build -o ./bin/qtun-darwin-amd64 ./main.go
#Windows amd64
GOOS=windows GOARCH=amd64 go build -o ./bin/qtun-windows-amd64.exe ./main.go

echo "DONE!!!"
