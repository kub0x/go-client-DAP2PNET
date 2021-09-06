package utils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func NewHTTPRequest(URL string, method string, data []byte) ([]byte, error) {
	var req *http.Request
	var err error
	if data != nil {
		req, err = http.NewRequest(method, URL, bytes.NewBuffer(data))
	} else {
		req, err = http.NewRequest(method, URL, nil)
	}
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return body, errors.New("unvalid status code")
	}

	return body, nil
}

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return body, errors.New(fmt.Sprintf("unvalid status code: %v", resp.StatusCode))
	}

	return body, nil

}
