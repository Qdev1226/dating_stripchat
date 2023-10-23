package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/kriskelly/dating-app-example/internal/dgraph"
)

func readSchema() []byte {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	path := filepath.Join(wd, "./internal/dgraph/dgraph.graphqls")
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	return b
}

func main() {
	log.Println("Setting up the Dgraph schema...")

	client := dgraph.NewClient()
	client.Connect()
	defer client.Close()
	op := &api.Operation{}
	op.Schema = string(readSchema())
	if err := client.Alter(context.Background(), op); err != nil {
		log.Fatal(err)
	}

	log.Println("Ran Alter Schema on DGraph")
}
