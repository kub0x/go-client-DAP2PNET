package rendezvous

import (
	"dap2pnet/client/models"
	"dap2pnet/client/utils"
	"encoding/json"
	"fmt"
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
	subReq := models.SubscribeRequest{
		Port: rendez.NodePort,
	}
	subBytes, err := json.Marshal(subReq)
	if err != nil {
		return err
	}

	_, err = utils.NewHTTPSRequest(rendez.PeerSuscribeEndpoint, "POST", subBytes, true)
	if err != nil {
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

	fmt.Println(peers)

	return nil
}
