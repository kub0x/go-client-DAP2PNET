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

func (bucks *Buckets) AddTriplet(triplet models.Triplet) {
	// calculate distances so hex -> int needed for XORing and sorting
	peerIntID, _ := big.NewInt(0).SetString(triplet.ID, 16)
	parentIntID, _ := big.NewInt(0).SetString(bucks.ParentID, 16)
	distance := parentIntID.Xor(parentIntID, peerIntID) // WTF Golang only 1 arg should do it - missing my cpp operators :/
	log2Dist := distance.BitLen() - 1                   // zero index so minus 1
	targetBucket := bucks.List[log2Dist]
	targetBucket.Add(triplet, *distance)
}

func (bucks *Buckets) PrintBuckets() {
	for i := 0; i < bucks.MaxLen; i++ {
		buck := bucks.List[i]
		for k, _ := range buck.Triplets {
			println("b" + strconv.Itoa(i) + ":" + k)
		}
	}
}
