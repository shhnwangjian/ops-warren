package lib

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"
)

/*
主题信息
CN=test.com

通用名称(CN)	test.com
-----------
签发者信息
CN=R3,O=Let's Encrypt,C=US

通用名称(CN)	R3
组织(O)	Let's Encrypt
国家(C)	US
------------
证书信息


 */

const (
	CertNotParse = "无法解析的证书"
)

type CertChain struct {
	CertChain []CertInfo
}

type CertInfo struct {
	Before,After time.Time
	Subject string
	Issuer string
	Sans []string
	OcspUrl []string
	CaUrl []string
	ExtKeyUsage []string
	Signature string
	IsCa bool
}

// certDecode 证书链解析
func (c *CertChain)  certDecode(content []byte,  cert *tls.Certificate) tls.Certificate {
	certDERBlock, restPEMBlock :=  pem.Decode(content)
	if len(restPEMBlock) != 0 {
		c.certDecode(restPEMBlock, cert)
	}
	cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
	return *cert
}

func (c *CertChain) ParseCertificate(content string)  error {
	// 获取证书信息 -----BEGIN CERTIFICATE-----   -----END CERTIFICATE-----
	// 这里返回的第二个值是证书中剩余的 block, 一般是rsa私钥 也就是 -----BEGIN RSA PRIVATE KEY 部分
	// 一般证书的有效期，组织信息等都在第一个部分里
	var cert tls.Certificate
	newCert := c.certDecode([]byte(content), &cert)
	if len(newCert.Certificate) ==0 {
		return errors.New(CertNotParse)
	}

	for i:=0; i< len(newCert.Certificate); i ++ {
		var ci CertInfo
		x509Cert, err := x509.ParseCertificate(newCert.Certificate[i])
		if err != nil {
			return err
		}
		ci.Subject = x509Cert.Subject.String()  // 主题信息
		ci.Issuer = x509Cert.Issuer.String() // 签发者信息
		ci.Before = x509Cert.NotBefore.UTC().Local() // 颁发日期
		ci.After = x509Cert.NotAfter.UTC().Local() // 截止日期
		ci.Sans = x509Cert.DNSNames // sans
		ci.OcspUrl = x509Cert.OCSPServer // ocsp_url
		ci.CaUrl = x509Cert.IssuingCertificateURL // caUrl
		ci.ExtKeyUsage = c.extKeyUsageChange(x509Cert.ExtKeyUsage) // extKeyUsage
		ci.Signature = x509Cert.SignatureAlgorithm.String() // 算法
		ci.IsCa = x509Cert.IsCA // 根证书
		c.CertChain = append(c.CertChain, ci)
	}
	return nil
}


func (c *CertChain) extKeyUsageChange(ml []x509.ExtKeyUsage) (extKeyUsage []string) {
	if len(ml) == 0{
		return
	}
	for _, line := range ml {
		switch line {
		case x509.ExtKeyUsageAny:
			extKeyUsage= append(extKeyUsage, "Any")
		case x509.ExtKeyUsageServerAuth:
			extKeyUsage= append(extKeyUsage, "ServerAuth")
		case x509.ExtKeyUsageClientAuth:
			extKeyUsage= append(extKeyUsage, "ClientAuth")
		case x509.ExtKeyUsageCodeSigning:
			extKeyUsage= append(extKeyUsage, "CodeSigning")
		case x509.ExtKeyUsageEmailProtection:
			extKeyUsage= append(extKeyUsage, "EmailProtection")
		case x509.ExtKeyUsageIPSECEndSystem:
			extKeyUsage= append(extKeyUsage, "IPSECEndSystem")
		case x509.ExtKeyUsageIPSECTunnel:
			extKeyUsage= append(extKeyUsage, "IPSECTunnel")
		case x509.ExtKeyUsageIPSECUser:
			extKeyUsage= append(extKeyUsage, "IPSECUser")
		case x509.ExtKeyUsageTimeStamping:
			extKeyUsage= append(extKeyUsage, "TimeStamping")
		case x509.ExtKeyUsageOCSPSigning:
			extKeyUsage= append(extKeyUsage, "OCSPSigning")
		case x509.ExtKeyUsageMicrosoftServerGatedCrypto:
			extKeyUsage= append(extKeyUsage, "MicrosoftServerGatedCrypto")
		case x509.ExtKeyUsageNetscapeServerGatedCrypto:
			extKeyUsage= append(extKeyUsage, "NetscapeServerGatedCrypto")
		case x509.ExtKeyUsageMicrosoftCommercialCodeSigning:
			extKeyUsage= append(extKeyUsage, "MicrosoftCommercialCodeSigning")
		case x509.ExtKeyUsageMicrosoftKernelCodeSigning:
			extKeyUsage= append(extKeyUsage, "MicrosoftKernelCodeSigning")
		}
	}
	return
}