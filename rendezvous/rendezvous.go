package rendezvous

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"dap2pnet/client/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Rendezvous struct {
	PeerSuscribeEndpoint string
	PeerExchangeEndpoint string
	NodePort             uint16
}

func NewRendezvous(port uint16) *Rendezvous {
	rendez := &Rendezvous{
		PeerSuscribeEndpoint: "https://rendezvous.dap2p.net/peers/subscribe",
		PeerExchangeEndpoint: "https://rendezvous.dap2p.net/peers/",
		NodePort:             port,
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
		MaxVersion:   tls.VersionTLS12, // TODO: ONLY FOR DEBUGGING PURPOSES (wireshark)
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{Transport: tr}

	subReq := models.SubscribeRequest{
		Port: rendez.NodePort,
	}

	subBytes, err := json.Marshal(subReq)
	if err != nil {
		return err
	}

	resp, err := httpClient.Post(rendez.PeerSuscribeEndpoint, "text/html", bytes.NewBuffer(subBytes))
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return err
	}
	println(resp.StatusCode)
	println(string(body))

	return nil
}

func (rendez *Rendezvous) TestPeerExchange() error {
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
		MaxVersion:   tls.VersionTLS12, // TODO: ONLY FOR DEBUGGING PURPOSES (wireshark)
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{Transport: tr}

	resp, err := httpClient.Get(rendez.PeerExchangeEndpoint)

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var peers models.PeerInfo
	err = json.Unmarshal(body, &peers)
	if err != nil {
		return err
	}

	println(resp.StatusCode)
	fmt.Println(peers)

	return nil
}
