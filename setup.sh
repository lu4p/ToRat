die() {
    echo
    echo "Setup failed: $1"
    exit 1
}

if [ "$EUID" -ne 0 ]; then
    echo "This installation script needs to be run as root."
    echo "sudo ./setup.sh"
    exit 1
fi

git clone https://github.com/lu4p/ToRat_client.git || die "ToRat_client could not be cloned"
git clone https://github.com/lu4p/ToRat_server.git || die "ToRat_Server could not be cloned"
git clone https://github.com/lu4p/genCert.git || die "genCert could not be cloned"

tee -a /etc/tor/torrc <<EOF
HiddenServiceDir /var/lib/tor/torat/
HiddenServiceVersion 3
HiddenServicePort 1337 127.0.0.1:1338
EOF

systemctl restart tor || die "Tor service could not be restarted"

hostname=$(sudo cat /var/lib/tor/torat/hostname) || die "Could not read Hostname from /var/lib/tor/torat/hostname"

cd ./genCert || die "Could not cd ./genCert"
go run genCert.go --ca --host $hostname || die "Could not generate tls certificate"
cp *.pem ../ToRat_server || die "Could not copy *.pem to ../ToRat_server"
cd .. || die "Could not cd .."
cert=$(cat ./ToRat_server/cert.pem) || die "Could not read ./ToRat_server/cert.pem"

conf=$(cat << EOF
package client

import "github.com/lu4p/ToRat_client/crypto"

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

rm ./ToRat_client/client/conf.go -f || die "Could not remove ./ToRat_client/client/conf.go"
tee -a ./ToRat_client/client/conf.go<<EOF
${conf}

EOF || die "Could not Write config to ./ToRat_client/client/conf.go"

