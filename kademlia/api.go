package kademlia

import (
	"dap2pnet/client/kademlia/buckets"
	"dap2pnet/client/models"
	"dap2pnet/client/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Kademlia struct {
	ttl              int
	buckets          *buckets.Buckets
	FindNodeEndpoint string
}

// Use Breadth first search having at most TTL levels

func NewKademliaAPI(bucks *buckets.Buckets) *Kademlia {
	return &Kademlia{
		ttl:              4,
		buckets:          bucks,
		FindNodeEndpoint: "/peers/key/",
	}
}

func (kad *Kademlia) findNodes(list *models.PeerInfo, triplet models.Triplet, key string, ttl int) error {

	if ttl >= kad.ttl {
		return nil
	}
	uri := fmt.Sprintf("http://%v:%v%v%v", triplet.IP, triplet.Port, kad.FindNodeEndpoint, key)
	fmt.Printf("Asking %v:%v for nearest nodes to key %v\n", triplet.IP, triplet.Port, key)
	body, err := utils.NewHTTPRequest(uri, "GET", nil)
	if err != nil {
		return err
	}

	var peers models.PeerInfo
	err = json.Unmarshal(body, &peers)
	if err != nil {
		return err
	}

	list.Triplets = append(list.Triplets, peers.Triplets...)
	for i := 0; i < len(peers.Triplets); i++ {
		err = kad.findNodes(list, peers.Triplets[i], key, ttl+1)
		if err != nil {
			return err
		}
	}

	return nil
}

func (kad *Kademlia) FindNodes(key string) (*models.PeerInfo, error) {
	finalPeers := &models.PeerInfo{}
	nearest := kad.buckets.NearestToKey(key)
	if nearest == nil {
		return nil, errors.New("no near nodes to key")
	}
	for j := 0; j < len(nearest.Triplets); j++ {
		err := kad.findNodes(finalPeers, nearest.Triplets[j], key, 1)
		if err != nil {
			log.Println(err)
		}
	}

	return finalPeers, nil
}

func (kad *Kademlia) StoreKey(key string) {

}

func (kad *Kademlia) RetrieveKey(key string) {

}

func (kad *Kademlia) NodeHealth() {

}
