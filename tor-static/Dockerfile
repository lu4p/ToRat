FROM ubuntu:rolling

RUN apt-get update && apt-get upgrade && apt-get install -y git golang tor

RUN apt-get -y install build-essential libtool autopoint

RUN mkdir -p /root/go/src/github.com/cretz/
RUN cd /root/go/src/github.com/cretz/ && git clone https://github.com/cretz/tor-static.git --recursive
RUN cd /root/go/src/github.com/cretz/tor-static && go run build.go build-all

RUN cd /root/go/src/github.com/cretz/tor-static && go run build.go package-libs

RUN cd /root/go/src/github.com/cretz/tor-static && rm libs.zip && mv libs.tar.gz /libs_linux.tar.gz

#get windows libs
RUN apt-get -y install wget
RUN cd /root/go/src/github.com/cretz/tor-static  && wget https://github.com/lu4p/tor-static/releases/download/2/tor-static-windows-amd64.zip
