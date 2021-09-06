package buckets

import (
	"crypto/x509"
	"dap2pnet/client/models"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"strconv"
)

type Buckets struct {
	ParentID string
	MaxLen   int
	List     map[int]*Bucket
	alpha    int
}

func NewBuckets() *Buckets {
	// load CN (ID) from certificate
	certBytes, err := os.ReadFile("./certs/id.pem")
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(certBytes)
	if block == nil || block.Type != "CERTIFICATE" {
		log.Fatal("failed to decode ca certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	bucks := &Buckets{
		ParentID: cert.Subject.CommonName,
		MaxLen:   256,
		List:     make(map[int]*Bucket), // allocs array of pointers for filling later
	}
	// instantiate buckets
	for i := 0; i < bucks.MaxLen; i++ {
		bucks.List[i] = NewBucket(i)
	}

	return bucks
}

func (bucks *Buckets) distanceToParent(keyID string) *big.Int {
	// calculate distances so hex -> int needed for XORing and sorting
	keyIntID, _ := big.NewInt(0).SetString(keyID, 16)
	parentIntID, _ := big.NewInt(0).SetString(bucks.ParentID, 16) // WTF Golang only 1 arg should do it - missing my cpp operators :/
	return parentIntID.Xor(parentIntID, keyIntID)
}

func (bucks *Buckets) AddTriplet(triplet models.Triplet) {
	distance := bucks.distanceToParent(triplet.ID)
	log2Dist := distance.BitLen() - 1 // zero index so minus 1
	//println(log2Dist)
	targetBucket := bucks.List[log2Dist]
	targetBucket.Add(triplet, *distance)
}

func (bucks *Buckets) NearestToKey(keyID string) *models.PeerInfo {
	distance := bucks.distanceToParent(keyID)
	log2Dist := distance.BitLen() - 1 // zero index so minus 1
	success := false
	peers := &models.PeerInfo{}
	for i := log2Dist; i >= 0; i-- {
		if len(bucks.List[i].Triplets) > 0 {
			for _, v := range bucks.List[i].Triplets {
				peers.Triplets = append(peers.Triplets, v)
			}
			success = true
			break
		}
	}
	if success {
		return peers
	}

	for i := log2Dist + 1; i < bucks.MaxLen; i++ {
		if len(bucks.List[i].Triplets) > 0 {
			for _, v := range bucks.List[i].Triplets {
				peers.Triplets = append(peers.Triplets, v)
			}
			success = true
			break
		}
	}
	if success {
		return peers
	}

	return nil
}

func (bucks *Buckets) PrintBuckets() {
	for i := 0; i < bucks.MaxLen; i++ {
		buck := bucks.List[i]
		for k, _ := range buck.Triplets {
			println("b" + strconv.Itoa(i) + ":" + k)
		}
	}
}
