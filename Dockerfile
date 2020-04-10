FROM lu4p/tor-static:latest AS torat-pre

COPY setup_docker.sh /
COPY . /ToRat
RUN cd /ToRat && go mod download -x

RUN mkdir -p /dist/server && mkdir -p /dist/client
RUN go get -v -u github.com/lu4p/genCert

FROM torat-pre AS torat
RUN /setup_docker.sh && rm /setup_docker.sh

# Build ToRat_server
RUN cd /ToRat/cmd/server && go build -o /dist/server/ToRat_server && cp banner.txt /dist/server/

RUN cd /go/pkg/mod/github.com/cretz/tor-static && tar -xf libs_linux.tar.gz
RUN cd /ToRat/cmd/client && go build --ldflags "-s -w" -tags "tor" -o /dist/client/client_linux && upx /dist/client/client_linux

RUN cd /go/pkg/mod/github.com/cretz/tor-static && unzip -o tor-static-windows-amd64.zip 

RUN cd /ToRat/cmd/client && env GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc-posix CXX=x86_64-w64-mingw32-g++-posix go build -tags "tor" --ldflags "-s -w -H windowsgui" -o /dist/client/client_windows.exe

RUN upx /dist/client/client_windows.exe --force
RUN mkdir /dist_ext

CMD (tor -f /torrc&) && cp /dist /dist_ext -rf && cd /dist/server/ && ./ToRat_server
