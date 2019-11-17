FROM lu4p/tor-static

RUN go get -v -u github.com/BurntSushi/xgb && \
	go get -v -u github.com/cretz/bine && \
	go get -v -u github.com/lu4p/cat && \
	go get -v -u github.com/lu4p/go-escalate && \
	go get -v -u github.com/lu4p/shred && \
	go get -v -u github.com/pierrre/archivefile/zip && \
	go get -v -u github.com/vova616/screenshot && \
    go get -v -u github.com/abiosoft/ishell && \
    go get -v -u github.com/fatih/color && \
    go get -v -u github.com/jinzhu/gorm && \
    go get -v -u github.com/jinzhu/gorm/dialects/sqlite && \
    go get -v -u github.com/dimiro1/banner/autoload && \
    go get -v -u golang.org/x/crypto/ed25519 && \
    go get -v -u golang.org/x/crypto/sha3 && \
    go get -v -u golang.org/x/net/proxy
RUN go get -v -u github.com/karalabe/xgo
RUN go get -v -u github.com/lu4p/genCert
RUN apt-get install -y upx && rm -rf /var/lib/apt/lists/*
