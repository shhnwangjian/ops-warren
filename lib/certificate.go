package lib

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"
)

func parseCertificate(content string) (err error, commonName, name string, before, after time.Time) {
	// 获取证书信息 -----BEGIN CERTIFICATE-----   -----END CERTIFICATE-----
	// 这里返回的第二个值是证书中剩余的 block, 一般是rsa私钥 也就是 -----BEGIN RSA PRIVATE KEY 部分
	// 一般证书的有效期，组织信息等都在第一个部分里
	certDERBlock, _ := pem.Decode([]byte(content))
	if certDERBlock == nil {
		return errors.New("无法解析的证书"), "", "", time.Now(), time.Now()
	}
	x509Cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		return err, "", "", time.Now(), time.Now()
	}
	commonName = x509Cert.Issuer.CommonName
	name = x509Cert.Subject.CommonName
	before = x509Cert.NotBefore
	after = x509Cert.NotAfter
	return
}
