package pki

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"dap2pnet/client/utils"
	"encoding/pem"
	"io/ioutil"
)

type PKI struct {
	IdentityEndpoint string
}

func NewPKI() *PKI {
	pki := &PKI{
		IdentityEndpoint: "https://pki.dap2p.net/pki/register",
	}

	return pki
}

func (pki *PKI) IssueIdentity() error {

	// Generate PEM encoded CSR and PEM encoded Private Key
	// Send CSR to PKI to get a signed certificate

	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) //csrng
	if err != nil {
		return err
	}

	csr := &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization:  []string{"organization"},
			Country:       []string{"country"},
			Province:      []string{"province"},
			Locality:      []string{"locality"},
			StreetAddress: []string{"address"},
		},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, csr, ecKey)
	if err != nil {
		return err
	}

	csrPEM := new(bytes.Buffer)
	err = pem.Encode(csrPEM, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})
	if err != nil {
		return err
	}

	privBytes, err := x509.MarshalECPrivateKey(ecKey)
	if err != nil {
		return err
	}

	privPEM := new(bytes.Buffer)
	err = pem.Encode(privPEM, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privBytes,
	})
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./certs/id.key", privPEM.Bytes(), 0600)
	if err != nil {
		return err
	}

	certChain, err := utils.NewHTTPSRequest(pki.IdentityEndpoint, "POST", csrPEM.Bytes(), false)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./certs/id.pem", certChain, 0600)
	if err != nil {
		return err
	}

	return nil

}
