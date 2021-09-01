package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
)

func NewHTTPSRequest(URL string, method string, data []byte, mutualTLS bool) ([]byte, error) {
	certPool := x509.NewCertPool()

	caBytes, err := ioutil.ReadFile("./certs/ca.pem")
	if err != nil {
		return nil, err
	}
	certPool.AppendCertsFromPEM(caBytes)

	// trust dap2pnet CA

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		RootCAs:    certPool,
	}

	if mutualTLS { // adapt request for mutual TLS
		tlsCertChain, err := tls.LoadX509KeyPair("./certs/id.pem", "./certs/id.key")
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCertChain}
	}

	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	var req *http.Request
	if data != nil {
		req, err = http.NewRequest(method, URL, bytes.NewBuffer(data))
	} else {
		req, err = http.NewRequest(method, URL, nil)
	}
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Transport: tr}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unvalid status code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil

}
