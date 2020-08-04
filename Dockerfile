FROM lu4p/tor-static:latest AS torat-pre
RUN go get -v -u github.com/lu4p/genCert
RUN GO111MODULE=on go get mvdan.cc/garble@c2079ac0a16b0a2546eb3de12bcb2fedbd74c454

RUN mkdir /ToRat
WORKDIR /ToRat
COPY go.mod .
COPY go.sum .

RUN go mod download -x

RUN mkdir -p /dist/server && mkdir -p /dist/client

FROM torat-pre AS torat
COPY . .
RUN ./setup_docker.sh

# Build ToRat_server
RUN cd ./cmd/server && go build -o /dist/server/ToRat_server && cp banner.txt /dist/server/

# ENV GOPRIVATE=github.com,gopkg.in,golang.org,google.golang.org

# Build Linux Client
RUN cd /go/pkg/mod/github.com/cretz/tor-static && tar -xf libs_linux.tar.gz
RUN cd ./cmd/client && garble -literals -seed=random build -tags "tor" -o /dist/client/client_linux && upx /dist/client/client_linux

# Build Windows Client
RUN cd /go/pkg/mod/github.com/cretz/tor-static && unzip -o tor-static-windows-amd64.zip 
RUN cd ./cmd/client && env GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ garble -literals -seed=random build -tags "tor" --ldflags "-H windowsgui" -o /dist/client/client_windows.exe
RUN upx /dist/client/client_windows.exe --force

CMD (tor -f /torrc&) && cp /dist /dist_ext -rf && ls /dist_ext && cd /dist_ext/dist/server/ && ./ToRat_server
