package main

import (
	"dap2pnet/client/pki"
	"dap2pnet/client/rendezvous"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	nodePort := uint16(1025 + rand.Intn(65535-1025+1))
	/*go func() {
		server.Run(nodePort) // TODO: Work with gin when testing is done
	}()*/

	pki := pki.NewPKI()

	err := pki.IssueIdentity()
	if err != nil {
		log.Fatal("couldnt get identity: " + err.Error())
	}

	rendez := rendezvous.NewRendezvous(nodePort)
	err = rendez.TestMutualTLS()
	if err != nil {
		log.Fatal("couldnt mutual tls: " + err.Error())
	}

	rendez.TestPeerExchange()
	time.Sleep(time.Minute * 1)
}
