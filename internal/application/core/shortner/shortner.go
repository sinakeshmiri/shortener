package shortner
import (
	"github.com/jxskiss/base62"
	"github.com/bwmarrin/snowflake"
)
// Shortner impliments the url-shortner interface
type Shortner struct {
	node *snowflake.Node
}

// NewShortner creates a new Shortner
func New(nodeName string) (*Shortner,error) {
	nodeID := djb2(nodeName)
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	return &Shortner{node: node},nil
}

// Short returns and id for the result of shortening the url
func (s Shortner)Short(url string) ([]byte, error) {
	id := s.node.Generate()
	return  base62.FormatInt(id.Int64()),nil
}


// hash function to map  node name to a unique int64
func djb2(inputString string) int64 {
	var  hashValue int64 = 0
	for _, char := range inputString {
		hashValue = (hashValue*31 + int64(char)) % 1024
	}
	return hashValue
}
