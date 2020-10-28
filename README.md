# qtun

A secure tunnel over quic.  

[![Travis](https://travis-ci.com/net-byte/qtun.svg?branch=master)](https://github.com/net-byte/qtun)
[![Go Report Card](https://goreportcard.com/badge/github.com/net-byte/qtun)](https://goreportcard.com/report/github.com/net-byte/qtun)
# Usage  

```
Usage of ./qtun:
 
  -from string
      From address (default "127.0.0.1:1987")
  -to string
      To address (default "127.0.0.1:1080")
  -ck string
        Client key (default "../certs/client.key")
  -cp string
        Client pem (default "../certs/client.pem")
  -sk string
        Server key (default "../certs/server.key")
  -sp string
        Server pem (default "../certs/server.pem")
  -t int
        Connection timeout of seconds (default 10)
  -s  Server mode   

```  

# Docker build  
```
docker build . -t qtun
```  

# Docker run client    
```
docker run -d --name qtun-client -p 1987:1987 qtun -from=:1987 -to=SERVER_IP:1988 -ck=/app/certs/client.key -cp=/app/certs/client.pem -sk=/app/certs/server.key -sp=/app/certs/server.pem
```

# Docker run server    
```
docker run -d --name qtun-server -p 1988:1988/udp qtun -s -from=:1988 -to=DST_IP:DST_PORT -ck=/app/certs/client.key -cp=/app/certs/client.pem -sk=/app/certs/server.key -sp=/app/certs/server.pem
```
