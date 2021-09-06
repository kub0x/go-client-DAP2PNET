package main

import (
	"dap2pnet/client/kademlia"
	"dap2pnet/client/kademlia/buckets"
	"dap2pnet/client/pki"
	"dap2pnet/client/rendezvous"
	"dap2pnet/client/server"
	"log"
	"math/rand"
	"time"
)

func exchange() {

	pki := pki.NewPKI()
	err := pki.IssueIdentity()
	if err != nil {
		log.Fatal("couldnt get identity: " + err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	nodePort := uint16(1025 + rand.Intn(65535-1025+1))
	rendez := rendezvous.NewRendezvous(nodePort)

	log.Printf("I listen on: %v\n", nodePort)

	go func(bucks *buckets.Buckets) {
		server.Run(nodePort, bucks) // TODO: Work with gin when testing is done
	}(rendez.Buckets)

	for true {
		err = rendez.TestMutualTLS()
		if err != nil {
			log.Fatal("couldnt mutual tls: " + err.Error())
		}
		rendez.TestPeerExchange()
		rendez.Buckets.PrintBuckets()
		time.Sleep(time.Minute * 2)
	}

}

func request() {
	pki := pki.NewPKI()
	err := pki.IssueIdentity()
	if err != nil {
		log.Fatal("couldnt get identity: " + err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	nodePort := uint16(1025 + rand.Intn(65535-1025+1))
	rendez := rendezvous.NewRendezvous(nodePort)

	log.Printf("I listen on: %v\n", nodePort)

	go func(bucks *buckets.Buckets) {
		server.Run(nodePort, bucks) // TODO: Work with gin when testing is done
	}(rendez.Buckets)

	for true {
		err = rendez.TestMutualTLS()
		if err != nil {
			log.Fatal("couldnt mutual tls: " + err.Error())
		}
		rendez.TestPeerExchange()

		kad := kademlia.NewKademliaAPI(rendez.Buckets)
		nearests, err := kad.FindNodes(rendez.Buckets.ParentID)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(nearests)

		for _, triplet := range nearests.Triplets {
			rendez.Buckets.AddTriplet(triplet)
		}
		rendez.Buckets.PrintBuckets()

		time.Sleep(time.Minute * 2)
	}

}

func main() {
	exchange()
}
