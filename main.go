package main

import (
	"dap2pnet/client/kademlia/buckets"
	"dap2pnet/client/pki"
	"dap2pnet/client/rendezvous"
	"dap2pnet/client/server"
	"log"
	"math/rand"
	"time"
)

func main() {
	pki := pki.NewPKI()
	err := pki.IssueIdentity()
	if err != nil {
		log.Fatal("couldnt get identity: " + err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	nodePort := uint16(1025 + rand.Intn(65535-1025+1))
	limit := time.Now().Add(time.Second * 45).UnixNano()
	rendez := rendezvous.NewRendezvous(nodePort)

	log.Printf("I listen on: %v\n", nodePort)

	go func(bucks *buckets.Buckets) {
		server.Run(4444, bucks) // TODO: Work with gin when testing is done
	}(rendez.Buckets)

	for time.Now().UnixNano() < limit {
		err := rendez.TestMutualTLS()
		if err != nil {
			log.Fatal("couldnt mutual tls: " + err.Error())
		}

		rendez.TestPeerExchange()

		// kad := kademlia.NewKademliaAPI(rendez.Buckets)
		// nearests, err := kad.FindNodes(rendez.Buckets.ParentID)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// fmt.Println(nearests)

		// for _, triplet := range nearests.Triplets {
		// 	rendez.Buckets.AddTriplet(triplet)
		// }

		time.Sleep(time.Millisecond * 250)
	}
	rendez.Buckets.PrintBuckets()
}
