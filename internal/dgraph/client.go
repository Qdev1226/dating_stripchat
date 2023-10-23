package dgraph

import (
	"flag"
	"log"
	"os"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"

	"google.golang.org/grpc"
)

var (
	dgraph = flag.String("d", os.Getenv("DGRAPH_HOST"), "Dgraph Alpha address")
)

// Client is our Dgraph API wrapper
type Client struct {
	conn *grpc.ClientConn
	*dgo.Dgraph
}

// NewClient generates an instance of our client
func NewClient() *Client {
	return &Client{}
}

// Connect to GRPC
func (c *Client) Connect() {
	flag.Parse()
	conn, err := grpc.Dial(*dgraph, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
	c.Dgraph = dgo.NewDgraphClient(api.NewDgraphClient(c.conn))
}

// Close the connection
func (c *Client) Close() error {
	return c.conn.Close()
}
