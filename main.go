package main

import (
	"dap2pnet/client/pki"
	"dap2pnet/client/rendezvous"
	"log"
)

func main() {
	pki := pki.NewPKI()

	err := pki.IssueIdentity()
	if err != nil {
		log.Fatal("couldnt get identity: " + err.Error())
	}

	rendez := rendezvous.NewRendezvous()
	err = rendez.TestMutualTLS()
	if err != nil {
		log.Fatal("couldnt mutual tls: " + err.Error())
	}
}
