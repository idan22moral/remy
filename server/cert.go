// reference: https://gist.github.com/annanay25/43e3846e21b30818d8dcd5f9987e852d

package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"time"
)

func generateCACertificate() (*x509.Certificate, *rsa.PrivateKey) {
	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("generating random key: %v", err)
	}

	caCertTmpl, err := createCertificateTemplate()
	if err != nil {
		log.Fatalf("creating cert template: %v", err)
	}

	// this cert will be the CA that we will use to sign the server cert
	caCertTmpl.IsCA = true
	// describe what the certificate will be used for
	caCertTmpl.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	caCertTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}

	caCertificate, _, err := createCertificate(caCertTmpl, caCertTmpl, &caKey.PublicKey, caKey)
	if err != nil {
		log.Fatalf("error creating cert: %v", err)
	}

	return caCertificate, caKey
}

func generateTLSCredentials(caCertificate *x509.Certificate, caKey *rsa.PrivateKey) (tls.Certificate, []byte) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("generating random key: %v", err)
	}

	certTmpl, err := createCertificateTemplate()
	if err != nil {
		log.Fatalf("creating cert template: %v", err)
	}
	certTmpl.KeyUsage = x509.KeyUsageDigitalSignature
	certTmpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}

	_, certPEM, err := createCertificate(certTmpl, caCertificate, &key.PublicKey, caKey)
	if err != nil {
		log.Fatalf("error creating cert: %v", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		log.Fatalf("invalid key pair: %v", err)
	}

	return tlsCert, keyPEM
}

func createCertificateTemplate() (*x509.Certificate, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.New("failed to generate serial number: " + err.Error())
	}

	tmpl := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Remy"}},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		BasicConstraintsValid: true,
	}
	return &tmpl, nil
}

func createCertificate(template, parent *x509.Certificate, pub interface{}, parentPriv interface{}) (
	cert *x509.Certificate, certPEM []byte, err error) {

	certDER, err := x509.CreateCertificate(rand.Reader, template, parent, pub, parentPriv)
	if err != nil {
		return
	}

	cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return
	}

	b := pem.Block{Type: "CERTIFICATE", Bytes: certDER}
	certPEM = pem.EncodeToMemory(&b)
	return
}
