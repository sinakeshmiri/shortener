package shortener

import (
	"errors"
	"net"
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/jxskiss/base62"
)

// shortener impliments the url-shortener interface
type shortener struct {
	node     *snowflake.Node
	urlschan chan string
}

// Newshortener creates a new shortener
func New(nodeName string, urlschan chan string) (*shortener, error) {
	var nodeID int64 
	nodeID,err := iplower(nodeName)
	if err != nil {
	nodeID = djb2(nodeName)
	}
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	return &shortener{urlschan: urlschan, node: node}, nil
}

// Short returns and id for the result of shortening the url
func (s shortener) Short() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	for {
		id := s.node.Generate()
		s.urlschan <- string(base62.FormatInt(id.Int64()))
	}
}

// hash function to map  node name to a unique int64
func djb2(inputString string) int64 {
	var hashValue int64 = 0
	for _, char := range inputString {
		hashValue = (hashValue*31 + int64(char)) % 1024
	}
	return hashValue
}

func iplower(inputString string) (int64,error) {
	ip := net.ParseIP(inputString)
	if ip == nil {
		return 0,errors.New("Invalid IP")
	}
	ipBytes := ip.To4()
	if ipBytes == nil {
		return 0,errors.New("IPv6 address not supported")
	}

	// Get the last 2 bytes of the IP address
	lastTwoBytes := uint16(ipBytes[2])<<8 | uint16(ipBytes[3])

	// Get the lower 10 bits
	lower10Bits := lastTwoBytes & 0x3FF
	return int64(lower10Bits),nil

}
