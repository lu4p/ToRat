FROM ubuntu:rolling

RUN apt-get update && apt-get upgrade && apt-get install -y git golang tor

RUN apt-get -y install build-essential libtool autopoint

RUN mkdir -p /root/go/src/github.com/cretz/
RUN cd /root/go/src/github.com/cretz/ && git clone https://github.com/cretz/tor-static.git --recursive
RUN cd /root/go/src/github.com/cretz/tor-static && go run build.go build-all