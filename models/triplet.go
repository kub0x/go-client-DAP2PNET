package models

type Triplet struct {
	ID   string `json name=id`
	IP   string `json name=ip`
	Port string `json name=port`
}

type PeerInfo struct {
	Triplets []Triplet `json name=triplets`
}

type SubscribeRequest struct {
	Port uint16 `json name=port`
}
