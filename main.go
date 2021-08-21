package main

import (
	"dap2pnet/client/client"
	"dap2pnet/client/pki"
	"log"
)

func main() {
	pki := pki.NewPKI()

	err := pki.IssueIdentity()
	if err != nil {
		log.Fatal("couldnt get identity: " + err.Error())
	}

	rendez := client.NewRendezvous()
	err = rendez.TestMutualTLS()
	if err != nil {
		log.Fatal("couldnt mutual tls: " + err.Error())
	}
}