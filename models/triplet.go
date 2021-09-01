package models

import "math/big"

type Triplet struct {
	ID       string `json name=id`
	IP       string `json name=ip`
	Port     string `json name=port`
	Distance big.Int
}

type PeerInfo struct {
	Triplets []Triplet `json name=triplets`
}

type SubscribeRequest struct {
	Port uint16 `json name=port`
}
