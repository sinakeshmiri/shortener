package shortner

import (
	"sync"

	"github.com/bwmarrin/snowflake"
	"github.com/jxskiss/base62"
)

// Shortner impliments the url-shortner interface
type Shortner struct {
	node *snowflake.Node
	urlschan chan string
}

// NewShortner creates a new Shortner
func New(nodeName string,urlschan chan string) (*Shortner,error) {
	
	nodeID := djb2(nodeName)
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	
	return &Shortner{urlschan : urlschan,node: node},nil
}

// Short returns and id for the result of shortening the url
func (s Shortner)Short() {
	wg:=sync.WaitGroup{}
	wg.Add(1)
	for {
		id := s.node.Generate()
		s.urlschan <- string(base62.FormatInt(id.Int64()))
	}
}


// hash function to map  node name to a unique int64
func djb2(inputString string) int64 {
	var  hashValue int64 = 0
	for _, char := range inputString {
		hashValue = (hashValue*31 + int64(char)) % 1024
	}
	return hashValue
}
