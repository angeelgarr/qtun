# qtun

A tcp proxy over quic.  

[![Travis](https://travis-ci.com/net-byte/qtun.svg?branch=master)](https://github.com/net-byte/qtun)
[![Go Report Card](https://goreportcard.com/badge/github.com/net-byte/qtun)](https://goreportcard.com/report/github.com/net-byte/qtun)
![image](https://img.shields.io/badge/License-MIT-orange)
![image](https://img.shields.io/badge/License-Anti--996-red)

# Usage  
## Cmd

```
Usage of ./qtun:
  -S    server mode
  -ck string
        client key file path (default "../certs/client.key")
  -cp string
        client pem file path (default "../certs/client.pem")
  -from string
        from address (default ":1987")
  -sk string
        server key file path (default "../certs/server.key")
  -sp string
        server pem file path (default "../certs/server.pem")
  -t int
        connection timeout in seconds (default 30)
  -to string
        to address (default ":1080")
```  

## Docker
### Build
```
docker build . -t qtun
```  

### Run client    
```
docker run -d --name qtun-client -p 1987:1987 qtun -from=:1987 -to=SERVER_IP:1988 -ck=/app/certs/client.key -cp=/app/certs/client.pem -sk=/app/certs/server.key -sp=/app/certs/server.pem
```

### Run server    
```
docker run -d --name qtun-server -p 1988:1988/udp qtun -s -from=:1988 -to=DST_IP:DST_PORT -ck=/app/certs/client.key -cp=/app/certs/client.pem -sk=/app/certs/server.key -sp=/app/certs/server.pem
```

## Setting on linux
It is recommended to increase the maximum buffer size by running:
```
sysctl -w net.core.rmem_max=2500000
```