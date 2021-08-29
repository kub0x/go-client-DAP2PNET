package pki

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
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

	certPool := x509.NewCertPool()

	caBytes, err := ioutil.ReadFile("./certs/ca.pem")
	if err != nil {
		return err
	}
	certPool.AppendCertsFromPEM(caBytes)

	// trust dap2pnet CA

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{Transport: tr}
	resp, err := httpClient.Post(pki.IdentityEndpoint, "text/html", bytes.NewBuffer(csrPEM.Bytes()))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("unvalid status code")
	}

	certChain, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./certs/id.pem", certChain, 0600)
	if err != nil {
		return err
	}

	return nil

}
