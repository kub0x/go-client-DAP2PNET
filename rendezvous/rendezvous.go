package rendezvous

import (
	"dap2pnet/client/kademlia/buckets"
	"dap2pnet/client/models"
	"dap2pnet/client/utils"
	"encoding/json"
	"log"
)

type Rendezvous struct {
	Buckets              *buckets.Buckets
	PeerSuscribeEndpoint string
	PeerExchangeEndpoint string
	NodePort             uint16
}

func NewRendezvous(port uint16) *Rendezvous {
	rendez := &Rendezvous{
		Buckets:              buckets.NewBuckets(),
		PeerSuscribeEndpoint: "https://rendezvous.dap2p.net/peers/subscribe",
		PeerExchangeEndpoint: "https://rendezvous.dap2p.net/peers/",
		NodePort:             port,
	}

	return rendez
}

func (rendez *Rendezvous) TestMutualTLS() error {
	subReq := models.SubscribeRequest{
		Port: rendez.NodePort,
	}
	subBytes, err := json.Marshal(subReq)
	if err != nil {
		return err
	}

	msg, err := utils.NewHTTPSRequest(rendez.PeerSuscribeEndpoint, "POST", subBytes, true)
	if err != nil {
		log.Println(string(msg))
		return err
	}

	return nil
}

func (rendez *Rendezvous) TestPeerExchange() error {
	body, err := utils.NewHTTPSRequest(rendez.PeerExchangeEndpoint, "GET", nil, true)
	if err != nil {
		return err
	}

	var peers models.PeerInfo
	err = json.Unmarshal(body, &peers)
	if err != nil {
		return err
	}

	for _, triplet := range peers.Triplets {
		rendez.Buckets.AddTriplet(triplet)
	}
	//println(len(peers.Triplets))
	//rendez.Buckets.PrintBuckets()

	return nil
}
