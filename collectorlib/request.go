package collectorlib

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/consul/consul/structs"
)

type Request struct {
	Root structs.CheckServiceNodes
}

func ParseRequest(input io.Reader) (Request, error) {
	decoder := json.NewDecoder(input)

	var root structs.CheckServiceNodes
	if err := decoder.Decode(&root); err != nil {
		return Request{}, err
	}

	return Request{
		Root: root,
	}, nil
}
