package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

type Rendezvous struct {
	PeerSuscribeEndpoint string
	PeerExchangeEndpoint string
}

func NewRendezvous() *Rendezvous {
	rendez := &Rendezvous{
		PeerSuscribeEndpoint: "https://rendezvous.dap2p.net/",
	}

	return rendez
}

func (rendez *Rendezvous) TestMutualTLS() error {

	certPool := x509.NewCertPool()

	caBytes, err := ioutil.ReadFile("./certs/ca.pem")
	if err != nil {
		return err
	}
	certPool.AppendCertsFromPEM(caBytes)

	tlsCertChain, err := tls.LoadX509KeyPair("./certs/id.pem", "./certs/id.key")
	if err != nil {
		return err
	}

	// trust dap2pnet CA

	tlsConfig := &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{tlsCertChain},
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Get(rendez.PeerSuscribeEndpoint)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	println(string(body))

	return nil
}
