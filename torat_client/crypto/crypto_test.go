package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"path/filepath"
	"reflect"
	"testing"
)

var testCert = `-----BEGIN CERTIFICATE-----
MIIFGDCCAwCgAwIBAgIQEbB9fUDrrBU3zfapiTdrDTANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMB4XDTIwMTExNDE1NDUzNVoXDTMwMTExMjE1NDUz
NVowEjEQMA4GA1UEChMHQWNtZSBDbzCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCC
AgoCggIBAM38Pz2WoiQxFcbTorbaHQfnvi7usempYysZpHVsvKo90xPwGx5HQL+A
v/u/OYtcUpE8P+orQ3Jxa0WK+9Aov29E5SjiB8rF005wbLWxo0tB+wSmIOnDoJGu
PlVV0b3m+sGPEoSKLBNLIp1HiiVmeTOwvJqf/3YaaZO54tBXL2Ysb9XdhI0YOv3t
/nBKpxYuLC5NMqQnQW5OTpMJRlSx6do/ijnWyQ075bozcDZ7ph3hfFSpOB9KkEWF
fzhtshQQw8Hzr1FBQFvOZH9F0ZFyx6/j4I8ndxd852NRrvtiEPPHkAA19cf+UBWx
A523cjPqWh1HfIapjzFd6TtXMBLHfJnD+EAyz5KPT27u0fcj/TthEmOplfAYoso9
CIGP67RmgRgX2mxxN24Tmuo4h2S+SklcuZm1szBArFe8GXX5Z+ljoNdrH5nKZn/F
uKjQRlw/Ojq+EAHP65fXb5Q1djI5zWzwkgjirRDjEaQRlGmJDXshpd7QaPMUm7Uz
7lUy6lvIWeaXr+3TWgyhrkunrBs4KWbhZAWenGo4zkO+u6rofy1B67V48sSNFI4G
/UWrQ1Kn0eFZysKlOAJsfCgvaTSNTVhRzpukP5HIFLJgTY7dzmo/x6M9lMJeqaS8
FG90pw8gVBgX2zfDEQd9OpuxXwVyO7xa/74pgajNdy/eU9muBVV3AgMBAAGjajBo
MA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggrBgEFBQcDATAPBgNVHRMBAf8E
BTADAQH/MB0GA1UdDgQWBBRvVDQ3HX8t6EA3hkM4+DIuh9tRxzARBgNVHREECjAI
ggYub25pb24wDQYJKoZIhvcNAQELBQADggIBAG7fSmTciLdtdGVAnNvQB9FyQwJ6
kLx+xGIk0XE63MxMsiYCe/zW9xmK41qfozofASo7pxa5RZ0PBUD67XhJ/YtK7EsD
cQQX4L5LY2xns5fzZ2aL/6V3/QlA+OywrvrsSW4QPidPOdmaCGgXRfwb/MAg1+ka
DK2EtrUYj6XhrY6RNNBY6KTlY4Wdivgd113q4fySIOCDM9Zb8jggFHIUawciI5Mh
3mMYWfJS9WLnELNTY+AaRsPh948uHP3wpp3YHeJtLKbqOD7PCPj98RN6urAmPp8V
ATQ7XPLcsBmYBEh54WpT/XJXaHoHsmNtdCA8NAyI4d2RjqCLou3CrOs7roe859EH
/IRNAjEdqx6L4BuM9IolB2ZfSGS3+dsxYh5b/tdOYOrLDYvWcny7OWkqUxRWBW+G
hPl10DvoFYYRRbhY9SFbkIuHMcVlXZaoRdtOEMUtN/nCjdh4AMJIrNYaQ2Da49L2
La0/AU8KnurLe8xV3hNO+mM27xARX/HvMAlQPSdFBsIiN0Rzbka3xSNLGQgptzSU
chOVvkFWj/2NpgfFljua9Qvy9UzACe8WcE1pU8ABqOHkvXuO4uagwmmeTnCLiWJQ
Eo+6BEZLgTK5VSk6+NJrHz03EFfVt8T16rC/SFflf5pg1x30qFFdeHw7rYkmbpW8
LuvpGsQoGfbMMIG2
-----END CERTIFICATE-----`

var pubKey = func() *rsa.PublicKey {
	block, _ := pem.Decode([]byte(testCert))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	return cert.PublicKey.(*rsa.PublicKey)
}()

func TestRandString(t *testing.T) {
	first := GenRandString()
	if len(first) != 16 {
		t.Fatal("wrong hostname length")
	}
	second := GenRandString()
	if first == second {
		t.Fatal("hostname not random")
	}
}

func TestGenHostname(t *testing.T) {
	hostname := genHostname(pubKey)
	dir := t.TempDir()

	path := filepath.Join(dir, "hostname")

	if err := encodeToFile(hostname, path); err != nil {
		t.Fatal(err)
	}

	readHostname, err := getEncodedFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(hostname, readHostname) {
		t.Fatal("generated and read hostname not equal")
	}
}
